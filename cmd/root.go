/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

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
		CrawlEfcCmd,
	)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
