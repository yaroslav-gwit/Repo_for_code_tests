package main

import (
	// "fmt"
	"io/ioutil"
	"log"
)

func main() {
	vmList()
}

func vmList() {
	folders, err := ioutil.ReadDir("/zroot/vm-encrypted/")
	if err != nil {
		log.Fatal(err)
	}

	var vm_list = []string{}

	for _, folder := range folders {
		vm_list = append(vm_list, folder.Name())
	}
	return vm_list
}
