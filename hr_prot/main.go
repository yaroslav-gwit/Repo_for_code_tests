package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	var vm_list = vmList()
	var live = "🟢"
	var not_live = "🔴"

	for _, vm := range vm_list {
		if vmLiveCheck(vm) {
			var vm_name = live + " " + vm
			fmt.Println(vm_name)
		} else {
			var vm_name = not_live + " " + vm
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
	vms, err := ioutil.ReadDir(bhyve_live_vms_folder)

	if err != nil {
		log.Fatal(err)
	}

	for index, vm := range vms {
		if vm.Name() == vmname {
			return true
		} else if vm.Name() != vmname && index != len(vms) {
			continue
		}
	}
	return false
}
