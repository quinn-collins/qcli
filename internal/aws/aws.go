package awsconsumer

import (
	"context"
	"fmt"
	"log"
	"os/user"
	"strings"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/quinn-collins/qcli/internal/app"
	"golang.org/x/term"
	"gopkg.in/ini.v1"
)

type Client struct {
	iam *iam.Client
	s3  *s3.Client
	sts *sts.Client
}

func New() *Client {
	app := app.New()
	fmt.Printf("%+v\n", app.Config)

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile(app.Config.AWSProfile),
		config.WithRegion(app.Config.AWSRegion),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("--- Region ---\nRegion: [ %s ]", cfg.Region)
	c := Client{
		iam: iam.NewFromConfig(cfg),
		s3:  s3.NewFromConfig(cfg),
		sts: sts.NewFromConfig(cfg),
	}

	return &c
}

func MFA() error {
	app := app.New()
	fmt.Printf("%+v", app.Config)

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile(app.Config.AWSProfile),
		config.WithRegion(app.Config.AWSRegion),
	)
	if err != nil {
		return fmt.Errorf("could not load aws config: %w", err)
	}

	iamClient := iam.NewFromConfig(cfg)
	mfaDevices, err := iamClient.ListMFADevices(context.TODO(), &iam.ListMFADevicesInput{})
	if err != nil {
		return fmt.Errorf("could not retrieve mfa devices: %w", err)
	}

	stsClient := sts.NewFromConfig(cfg)

	fmt.Print("Enter MFA token: ")
	byteMFAToken, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("could not read mfa token: %w", err)
	}

	mfaToken := strings.TrimSpace(string(byteMFAToken))

	session, err := stsClient.GetSessionToken(context.TODO(), &sts.GetSessionTokenInput{
		SerialNumber: mfaDevices.MFADevices[0].SerialNumber,
		TokenCode:    aws.String(mfaToken),
	})
	if err != nil {
		return fmt.Errorf("could not get session token: %w", err)
	}

	user, err := user.Current()
	if err != nil {
		return fmt.Errorf("could not get current user: %w", err)
	}

	filePath := user.HomeDir + "/.aws/credentials"

	file, err := ini.Load(filePath)
	if err != nil {
		return fmt.Errorf("could not open credentials file: %s: %w", filePath, err)
	}

	section, err := file.NewSection(app.Config.AWSTargetProfile)
	if err != nil {
		return fmt.Errorf("could not create new section: %s: %w", app.Config.AWSTargetProfile, err)
	}

	expirationTime := session.Credentials.Expiration

	_, err = section.NewKey("aws_access_key_id", *session.Credentials.AccessKeyId)
	if err != nil {
		log.Println(err)
	}
	_, err = section.NewKey("aws_secret_access_key", *session.Credentials.SecretAccessKey)
	if err != nil {
		log.Println(err)
	}
	_, err = section.NewKey("aws_session_token", *session.Credentials.SessionToken)
	if err != nil {
		log.Println(err)
	}
	_, err = section.NewKey("expiration", expirationTime.String())
	if err != nil {
		log.Println(err)
	}

	err = file.SaveTo(filePath)
	if err != nil {
		return fmt.Errorf("could not save to credentials file: %w", err)
	}

	fmt.Printf("Profile: [ %s ] temporary credentials will expire in %s at %s", app.Config.AWSTargetProfile, expirationTime.Sub(time.Now()).Round(60*time.Second).String(), expirationTime.Format(time.RFC3339))

	return nil
}

func (c *Client) Me() *sts.GetCallerIdentityOutput {
	result, err := c.sts.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func (c *Client) ListBuckets() *s3.ListBucketsOutput {
	result, err := c.s3.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func (c *Client) ListObjects(bucket string) *s3.ListObjectsV2Output {
	result, err := c.s3.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		log.Fatal(err)
	}

	return result
}
