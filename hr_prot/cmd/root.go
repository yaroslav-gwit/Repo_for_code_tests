package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hoster",
	Short: "HosterRed is a fast, CLI-based Bhyve config and VM manager written in Go",
	// This can be a multiline string!
	Long: `HosterRed is a fast, CLI-based Bhyve config and VM manager written in Go`,

	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
