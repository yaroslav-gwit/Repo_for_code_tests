package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
}

func vmList(plain ...bool) []string {
	var folder_to_scan = "/zroot/vm-encrypted/"
	folders, err := ioutil.ReadDir(folder_to_scan)
	if err != nil {
		log.Fatal(err)
	}

	var vm_list = []string{}

	for _, folder := range folders {
		var vm_folder = folder_to_scan + folder.Name()
		var vm_folder_files, _ = ioutil.ReadDir(vm_folder)

		for _, file := range vm_folder_files {
			if file.Name() == "vm.config" || file.Name() == "vm.conf" {
				vm_list = append(vm_list, folder.Name())
			}
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
