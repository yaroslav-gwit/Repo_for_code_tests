package main

import (
	// "fmt"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	fmt.Println(vmList())
}

func vmList() []string {
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
