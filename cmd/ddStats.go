/*
Copyright © 2024 Quinn Collins <collinsqui@gmail.com>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/spf13/cobra"
)

// listBucketsCmd represents the listBuckets command
var ddStats = &cobra.Command{
	Use:   "dd-stats",
	Short: "do some dd stuff",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("dd stats called")

		_, err := statsd.New("127.0.0.1:8125")
		if err != nil {
			log.Fatal(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(listBucketsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listBucketsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listBucketsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
