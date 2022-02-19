package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print out the version number of HosterRed",
	Long:  `Print out the version number and/or development release of the HosterRed`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("HosterRed Release 0.01-alpha")
	},
}
