package cmd

import (
	"fmt"
	"os"

	"github.com/shirou/gopsutil/mem"
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
		// Hostname
		hostname, _ := os.Hostname()
		fmt.Println("Hostname: " + hostname)
		// RAM
		vmStat, _ := mem.VirtualMemory()
		fmt.Println("Free RAM:", vmStat.Available)
		fmt.Println("Overall RAM:", vmStat.HighTotal)
	},
}
