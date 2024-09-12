/*
Copyright © 2024 Quinn Collins <collinsqui@gmail.com>
*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "qcli",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println(cmd.Use)
	// },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Set defaults
	viper.SetDefault("aws-profile", "default")
	viper.SetDefault("aws-target-profile", "default")
	viper.SetDefault("aws-region", "us-east-1")

	// Find home directory
	home, err := os.UserHomeDir()
	if err != nil {
		log.Println(fmt.Errorf("could not find home directory: %w", err))
	}

	// Read and parse configuration file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(home + "/.config/qcli")
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			// Config file not found
			// I plan to leave this empty later
			log.Println(fmt.Errorf("config file not found: %w", err))
		}
	}

	// Read and parse environment variables
	err = viper.BindEnv("aws-region", "QCLI_AWS_REGION")
	err = viper.BindEnv("aws-profile", "QCLI_AWS_PROFILE")
	err = viper.BindEnv("aws-target-profile", "QCLI_AWS_TARGET_PROFILE")
	if err != nil {
		fmt.Println(fmt.Errorf("could not bind environment variable: %w", err))
	}
}
