package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	fmt.Println(vmList())
}

func vmList() []string {
	var folder_to_scan = "/zroot/vm-encrypted/"
	folders, err := ioutil.ReadDir(folder_to_scan)
	if err != nil {
		log.Fatal(err)
	}

	var vm_list = []string{}

	for _, folder := range folders {
		_folder := folder_to_scan + folder.Name()
		_files, _ := ioutil.ReadDir(_folder)
		for _, _file := range _files {
			if _file.Name() == "vm.config" || _file.Name() == "vm.conf" {
				vm_list = append(vm_list, folder.Name())
			}
		}
	}
	return vm_list
}
