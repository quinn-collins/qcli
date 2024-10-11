/*
Copyright © 2024 Quinn Collins <collinsqui@gmail.com>
*/
package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc"
	"github.com/quinn-collins/qcli/internal/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ssoCmd represents the sso command
var ssoCmd = &cobra.Command{
	Use:   "sso",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("aws-profile", cmd.Flags().Lookup("aws-profile"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.TODO()
		app.New()

		var ssoStartURL, ssoRegion string

		awsCfg, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(app.NewConfig().AWSProfile))
		if err != nil {
			panic(err)
		}

		for _, source := range awsCfg.ConfigSources {
			switch s := source.(type) {
			case config.SharedConfig:
				if s.SSOSession != nil {
					ssoStartURL = s.SSOSession.SSOStartURL
					ssoRegion = s.SSOSession.SSORegion
				}
			}
		}

		ssoOIDCClient := ssooidc.NewFromConfig(awsCfg, func(o *ssooidc.Options) {
			o.Region = ssoRegion
		})
		registerClientOutput, err := ssoOIDCClient.RegisterClient(ctx, &ssooidc.RegisterClientInput{
			ClientName: aws.String("qcli"),
			ClientType: aws.String("public"),
		})
		if err != nil {
			panic(err)
		}

		startDeviceAuthorizationOutput, err := ssoOIDCClient.StartDeviceAuthorization(ctx, &ssooidc.StartDeviceAuthorizationInput{
			ClientId:     registerClientOutput.ClientId,
			ClientSecret: registerClientOutput.ClientSecret,
			StartUrl:     aws.String(ssoStartURL),
		})
		if err != nil {
			panic(err)
		}

		open(*startDeviceAuthorizationOutput.VerificationUriComplete)

		var createTokenOutput *ssooidc.CreateTokenOutput
		// Refactor this to get createTokenOutput on a channel and timeout after 30 seconds instead of polling every 500ms
		sendMessage := true
		for {
			createTokenOutput, err = ssoOIDCClient.CreateToken(ctx, &ssooidc.CreateTokenInput{
				ClientId:     registerClientOutput.ClientId,
				ClientSecret: registerClientOutput.ClientSecret,
				GrantType:    aws.String("urn:ietf:params:oauth:grant-type:device_code"),
				DeviceCode:   startDeviceAuthorizationOutput.DeviceCode,
			})
			if err != nil {
				var re *awshttp.ResponseError
				if errors.As(err, &re) {
					isPending := strings.Contains(re.Unwrap().Error(), "AuthorizationPendingException")
					if sendMessage {
						sendMessage = false
						fmt.Println("Attempting to open SSO authorization page in your browser.")
						fmt.Println("\n", *startDeviceAuthorizationOutput.VerificationUri)
						fmt.Println("\nAlternatively, you may open provided URL on any device and provide the following code.")
						fmt.Println(*startDeviceAuthorizationOutput.UserCode)
					}
					if isPending {
						time.Sleep(500 * time.Millisecond)
						continue
					} else {
						panic(err)
					}
				} else {
					panic(err)
				}
			}
			break
		}

		expiresIn := time.Duration((time.Duration(createTokenOutput.ExpiresIn) * time.Second))
		expiresAt := time.Now().Add(expiresIn)
		clientSecretExpiresAt := time.Now().Add(time.Duration(registerClientOutput.ClientSecretExpiresAt) * time.Millisecond)

		accessTokenData := struct {
			AccessToken           string
			ExpiresAt             string
			ExpiresIn             string
			StartURL              string
			Region                string
			ClientID              string
			ClientSecret          string
			ClientSecretExpiresAt string
			RefreshToken          string
		}{
			AccessToken:           *createTokenOutput.AccessToken,
			ExpiresAt:             expiresAt.String(),
			ExpiresIn:             expiresIn.String(),
			StartURL:              ssoStartURL,
			Region:                ssoRegion,
			ClientID:              *registerClientOutput.ClientId,
			ClientSecret:          *registerClientOutput.ClientSecret,
			ClientSecretExpiresAt: clientSecretExpiresAt.String(),
		}

		accessTokenJSON, err := json.Marshal(accessTokenData)
		if err != nil {
			panic(err)
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		writePath := path.Join(homeDir, ".aws/sso/cache", "qcli-access-token.json")

		err = os.WriteFile(writePath, accessTokenJSON, 0666)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(ssoCmd)

	ssoCmd.Flags().StringP("aws-profile", "p", "", "AWS identity profile.")
}

func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
