/*
Copyright © 2024 Quinn Collins <collinsqui@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/quinn-collins/qcli/internal/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// meCmd represents the me command
var meCmd = &cobra.Command{
	Use:   "me",
	Short: "Get AWS caller details",
	Run: func(cmd *cobra.Command, args []string) {
		app := app.New()
		fmt.Printf("%+v\n", app.Config)
		// app := Application{
		// 	github: githubconsumer.New(),
		// 	aws:    awsconsumer.New(aws.String(awsRegion), aws.String(awsTargetProfile)),
		// }
		//
		// result := app.aws.Me()
		//
		// fmt.Printf("--- User Details ---\nArn: [ %s ]\nUserID: [ %s ]\nAccount: [ %s ]\n", *result.Arn, *result.UserId, *result.Account)
	},
}

func init() {
	rootCmd.AddCommand(meCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// meCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// meCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// Read and parse flags and bind them to viper
	meCmd.Flags().StringP("aws-profile", "p", "", "AWS identity profile.")
	err := viper.BindPFlag("aws-profile", meCmd.Flags().Lookup("aws-profile"))
	meCmd.Flags().StringP("aws-target-profile", "t", "", "Use a target aws-profile")
	err = viper.BindPFlag("aws-target-profile", meCmd.Flags().Lookup("aws-target-profile"))
	meCmd.Flags().StringP("aws-region", "r", "", "Select an AWS region. E.g. us-west-2")
	err = viper.BindPFlag("aws-region", meCmd.Flags().Lookup("aws-region"))

	fmt.Println(err)
}
