package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/alexstory/otters/cmd/ottr/new"
	"github.com/alexstory/otters/cmd/ottr/version"
)

var rootCmd = &cobra.Command{
	Use:   "ottr",
	Short: "ottr helps you manage your otters web app.",
	Long: `ottr has cli scaffolding and migration tools for your app.
visit github.com/alexstory/otters.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hi from otters. squeak squeak.")
	},
}

func init() {
	rootCmd.AddCommand(version.VersionCmd)
	rootCmd.AddCommand(new.NewCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
