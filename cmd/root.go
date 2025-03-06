/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/npvu1510/crawl-en-vocab/cmd/crawl"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "server",
	Short:        "start server",
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.AddCommand(
		crawl.CrawlEfcCmd,
	)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
