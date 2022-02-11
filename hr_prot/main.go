package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/jedib0t/go-pretty/v6/table"
	"gopkg.in/yaml.v2"
)

type datasetZfs struct {
	name       string
	mount_path string
	zfs_path   string
	encrypted  bool
}

//VM status icons
const vm_is_live = "🟢"
const vm_is_not_live = "🔴"
const vm_is_encrypted = "🔒"

func main() {
	var vm_list = vmList()

	var outputTable = table.NewWriter()
	outputTable.SetOutputMirror(os.Stdout)

	outputTable.AppendHeader(table.Row{"#", "vm name", "status"})

	for index, vm := range vm_list {
		var vm_status = ""
		if vmLiveCheck(vm) {
			vm_status = vm_is_live + vm_is_encrypted
		} else {
			vm_status = vm_is_not_live
		}
		outputTable.AppendRow([]interface{}{index + 1, vm, vm_status})
		outputTable.AppendSeparator()
	}

	var total_number_of_vms = strconv.Itoa(len(vm_list))
	outputTable.AppendFooter(table.Row{"", "total vms: " + total_number_of_vms})

	datasetZfsList()

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

func datasetZfsList() string {
	var conf_datasets_file, conf_datasets_error = os.ReadFile("conf_datasets.yaml")

	if conf_datasets_error != nil {
		panic(conf_datasets_error)
	}

	var datasetZfs_var = datasetZfs{}
	err := yaml.Unmarshal([]byte(conf_datasets_file), &datasetZfs_var)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Print(datasetZfs_var.name)

	var conf_datasets_data_output = string(conf_datasets_file)

	return conf_datasets_data_output
}
