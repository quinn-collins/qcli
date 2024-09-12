/*
Copyright © 2024 Quinn Collins <collinsqui@gmail.com>
*/
package cmd

import (
	"fmt"

	awsconsumer "github.com/quinn-collins/qcli/internal/aws"
	"github.com/quinn-collins/qcli/internal/tui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listBucketsCmd represents the listBuckets command
var listBucketsCmd = &cobra.Command{
	Use:   "list-buckets",
	Short: "Fetch buckets in AWS",
	Run: func(cmd *cobra.Command, args []string) {
		awsClient := awsconsumer.New()
		results := awsClient.ListBuckets()

		fmt.Printf("--- Buckets ---\n")
		for _, b := range results.Buckets {
			fmt.Printf("[ %s ] [ %s ]\n", *b.Name, *b.CreationDate)

		}
		tui.PrintTable(results.Buckets)
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		// https://github.com/spf13/viper/issues/233 This is why binding of flags is here
		viper.BindPFlag("aws-profile", cmd.Flags().Lookup("aws-profile"))
		viper.BindPFlag("aws-region", cmd.Flags().Lookup("aws-region"))
	},
}

func init() {
	rootCmd.AddCommand(listBucketsCmd)

	// Read and parse flags and bind them to viper
	listBucketsCmd.Flags().StringP("aws-profile", "p", "", "AWS identity profile.")
	listBucketsCmd.Flags().StringP("aws-region", "r", "", "Select an AWS region. E.g. us-west-2")
}
