package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/facette/natsort"
	"github.com/jedib0t/go-pretty/v6/table"
	"gopkg.in/yaml.v2"
)

type datasetsListStruct struct {
	Datasets []struct {
		Name       string `yaml:"name"`
		Mount_path string `yaml:"mount_path"`
		Zfs_path   string `yaml:"zfs_path"`
		Encrypted  bool   `yaml:"encrypted"`
		Type       string `yaml:"type"`
	}
}

type vmListStruct struct {
	vmName    []string
	vmDataset []string
}

type vmStatusCheckStruct struct {
	vmLive        bool
	vmEncrypted   bool
	vmStatusIcons string
}

func main() {
	var vm_list = vmList()

	var outputTable = table.NewWriter()
	outputTable.SetOutputMirror(os.Stdout)
	outputTable.AppendHeader(table.Row{"#", "vm name", "status", "dataset"})

	var vm_status = ""
	var vm_dataset = ""

	for index, vm := range vm_list.vmName {
		vm_status = vmStatusCheck(vm).vmStatusIcons
		vm_dataset = vm_list.vmDataset[index]
		outputTable.AppendRow([]interface{}{index + 1, vm, vm_status, vm_dataset})
		outputTable.AppendSeparator()
	}

	var total_number_of_vms = strconv.Itoa(len(vm_list.vmName))
	outputTable.AppendFooter(table.Row{"", "total vms: " + total_number_of_vms})

	fmt.Println("VM List")
	outputTable.SetStyle(table.StyleLight)
	outputTable.Render()
}

func vmList(plain ...bool) vmListStruct {
	var datasetsList_var = datasetsList()
	var folder_to_scan string
	var vm_list = vmListStruct{}

	for _, dataset := range datasetsList_var.Datasets {
		folder_to_scan = dataset.Mount_path
		folders, err := ioutil.ReadDir(folder_to_scan)
		if err != nil {
			log.Fatal(err)
		}

		for _, folder := range folders {
			var vm_folder_full_path = folder_to_scan + folder.Name()
			var vm_folder_name = folder.Name()

			var _, file_exists_error = os.Stat(vm_folder_full_path + "/vm.config")
			if file_exists_error == nil {
				vm_list.vmName = append(vm_list.vmName, vm_folder_name)
				// vm_list.vmDataset = append(vm_list.vmDataset, dataset.Name)
			}

			var _, new_config_file_exists_error = os.Stat(vm_folder_full_path + "/vm.conf")
			if new_config_file_exists_error == nil {
				vm_list.vmName = append(vm_list.vmName, vm_folder_name)
				// vm_list.vmDataset = append(vm_list.vmDataset, dataset.Name)
			}
		}
	}
	natsort.Sort(vm_list.vmName)

	for _, vm := range vm_list.vmName {
		for _, dataset := range datasetsList_var.Datasets {
			folder_to_scan = dataset.Mount_path
			var _, vm_in_dataset_error = os.Stat(folder_to_scan + vm)
			if vm_in_dataset_error == nil {
				vm_list.vmDataset = append(vm_list.vmDataset, dataset.Name)
			}
		}
	}

	return vm_list
}

func vmStatusCheck(vmname string) vmStatusCheckStruct {
	//VM status icons
	const vm_is_live = "🟢"
	const vm_is_not_live = "🔴"
	const vm_is_encrypted = "🔒"

	var vmStatusCheckStruct_var = vmStatusCheckStruct{}
	var vmStatusIcons = ""

	//VM live check
	var bhyve_live_vms_folder = "/dev/vmm/"
	if _, err := os.Stat(bhyve_live_vms_folder + vmname); err == nil {
		vmStatusCheckStruct_var.vmLive = true
		vmStatusIcons = vm_is_live
	} else {
		vmStatusCheckStruct_var.vmLive = false
		vmStatusIcons = vm_is_not_live
	}

	//VM encryption check
	var datasetsList_var = datasetsList()
	for index, dataset := range datasetsList_var.Datasets {
		var _, err = os.Stat(dataset.Mount_path + vmname)
		if err == nil {
			if dataset.Encrypted {
				vmStatusIcons = vmStatusIcons + vm_is_encrypted
				vmStatusCheckStruct_var.vmEncrypted = true
			}
		} else if err != nil && index != len(datasetsList_var.Datasets) {
			continue
		} else {
			vmStatusCheckStruct_var.vmEncrypted = false
			break
		}
	}

	vmStatusCheckStruct_var.vmStatusIcons = vmStatusIcons
	return vmStatusCheckStruct_var
}

func datasetsList() datasetsListStruct {
	var conf_datasets_file, conf_datasets_error = os.ReadFile("conf_datasets.yaml")

	if conf_datasets_error != nil {
		panic(conf_datasets_error)
	}

	var datasetsList_var = datasetsListStruct{}

	err := yaml.Unmarshal([]byte(conf_datasets_file), &datasetsList_var)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return datasetsList_var
}
