/*
Copyright © 2024 Quinn Collins <collinsqui@gmail.com>
*/
package cmd

import (
	"fmt"
	"log"

	awsconsumer "github.com/quinn-collins/qcli/internal/aws"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// mfaCmd represents the mfa command
var mfaCmd = &cobra.Command{
	Use:   "mfa",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := awsconsumer.MFA()
		if err != nil {
			log.Println(fmt.Errorf("failed to authenticate: %w", err))
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("aws-profile", cmd.Flags().Lookup("aws-profile"))
		viper.BindPFlag("aws-target-profile", cmd.Flags().Lookup("aws-target-profile"))
		viper.BindPFlag("aws-region", cmd.Flags().Lookup("aws-region"))
	},
}

func init() {
	rootCmd.AddCommand(mfaCmd)

	// Read and parse flags and bind them to viper
	mfaCmd.Flags().StringP("aws-profile", "p", "", "AWS identity profile.")
	mfaCmd.Flags().StringP("aws-target-profile", "t", "", "Use a target aws-profile")
	mfaCmd.Flags().StringP("aws-region", "r", "", "Select an AWS region. E.g. us-west-2")
}
