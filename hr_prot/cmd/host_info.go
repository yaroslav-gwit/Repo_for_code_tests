package cmd

import (
	"fmt"
	"os"

	"github.com/shirou/gopsutil/host"
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
		hostInfo()
	},
}

func hostInfo() {
	fmt.Println("This will be a host info!")
	// Hostname
	hostname, _ := os.Hostname()
	fmt.Println("Hostname: " + hostname)
	// RAM
	vmStat, _ := mem.VirtualMemory()
	freeRam := vmStat.Available / 1024 / 1024 / 1024
	totalRam := vmStat.Total / 1024 / 1024 / 1024
	fmt.Printf("RAM: %dGB/%dGB\n", freeRam, totalRam)
	//Uptime
	hostUptime := hostUptimeFunc()
	fmt.Println(hostUptime)
}

func hostUptimeFunc() string {
	hostUptime_, _ := host.Info()
	hostUptime := hostUptime_.Uptime
	// seconds := hostUptime % 60
	minutes := (hostUptime / 60) % 60
	hours := (hostUptime / 60 / 60) % 24
	days := hostUptime / 60 / 60 / 24
	// fmt.Printf("Uptime: %d days %d hours %d minutes", days, hours, minutes)
	return fmt.Sprintf("Uptime: %d days %d hours %d minutes", days, hours, minutes)
}
