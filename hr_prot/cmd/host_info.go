package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(hostinfoCmd)
}

var hostinfoCmd = &cobra.Command{
	Use:   "hostinfo",
	Short: "Show the host information",
	Long:  `Show the host information in the form of table or json output`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("This will be a host info!")
		hostname, _ := os.Hostname()
		fmt.Println("Hostname: " + hostname)
	},
}
