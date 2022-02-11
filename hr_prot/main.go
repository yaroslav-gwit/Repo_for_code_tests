package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

//VM status icons
const vm_is_live = "🟢"
const vm_is_not_live = "🔴"
const vm_is_encrypted = "🔒"

func main() {
	var vm_list = vmList()

	for _, vm := range vm_list {
		if vmLiveCheck(vm) {
			var vm_name = vm + " " + vm_is_live + vm_is_encrypted
			fmt.Println(vm_name)
		} else {
			var vm_name = vm + " " + vm_is_not_live
			fmt.Println(vm_name)
		}
	}

	//TABLE
	var outputTable = table.NewWriter()
	outputTable.SetOutputMirror(os.Stdout)

	outputTable.AppendHeader(table.Row{"#", "VM Name", "Status"})

	for index, vm := range vm_list {
		if vmLiveCheck(vm) {
			var vm_status = vm_is_live + vm_is_encrypted
			outputTable.AppendRow([]interface{}{index + 1, vm, vm_status})
		}
	}

	outputTable.AppendFooter(table.Row{"", "", "Total", 10000})

	outputTable.SetStyle(table.StyleLight)
	outputTable.Render()
}

func vmList(plain ...bool) []string {
	var folder_to_scan = "/zroot/vm-encrypted/"
	folders, err := ioutil.ReadDir(folder_to_scan)
	if err != nil {
		log.Fatal(err)
	}

	var vm_list = []string{}

	for _, folder := range folders {
		var vm_folder_full_path = folder_to_scan + folder.Name()
		var vm_folder_name = folder.Name()
		var _, file_exists_error = os.Stat(vm_folder_full_path + "/vm.config")

		if file_exists_error == nil {
			vm_list = append(vm_list, vm_folder_name)
		}
	}

	return vm_list
}

func vmLiveCheck(vmname string) bool {
	var bhyve_live_vms_folder = "/dev/vmm/"
	if _, err := os.Stat(bhyve_live_vms_folder + vmname); err == nil {
		return true
	} else {
		return false
	}
}
