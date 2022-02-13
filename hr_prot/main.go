package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/facette/natsort"
	"github.com/jedib0t/go-pretty/v6/table"
	"gopkg.in/yaml.v2"
)

type storageStruct struct {
	DiskSize  int
	DiskFreeG int
	DiskFreeP int
}

type vmConfigStruct struct {
	Cpus        int    `yaml:"cpus"`
	Ram         int    `yaml:"ram"`
	VncPort     int    `yaml:"vnc_port"`
	VncPassword string `yaml:"vnc_password"`
	IpAddress   string `yaml:"ip_address"`
	OsType      string `yaml:"os_type"`
	ParentHost  string `yaml:"parent_host"`
	Networks    []struct {
		InterfaceName      string `yaml:"interface_name"`
		InterfaceIpAddress string `yaml:"interface_ip_address"`
	}
	Storage []struct {
		DiskName     string `yaml:"disk_name"`
		DiskLocation string `yaml:"disk_location"`
	}
}

func VmConfig(vmname string) vmConfigStruct {
	var conf_vm_file, conf_vm_error = os.ReadFile("conf_vm.yaml")

	if conf_vm_error != nil {
		panic(conf_vm_error)
	}

	var vmConfigStruct_var vmConfigStruct

	err := yaml.Unmarshal([]byte(conf_vm_file), &vmConfigStruct_var)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return vmConfigStruct_var
}

func main() {
	var vm_list = vmList()

	var outputTable = table.NewWriter()
	outputTable.SetOutputMirror(os.Stdout)
	outputTable.AppendHeader(table.Row{"#", "vm name", "status", "dataset", "resources", "vnc", "networks", "storage", "misc"})

	var vm_status string
	var vm_dataset string
	var vm_resources string
	var vm_vnc string
	var vm_networks string
	var vm_os_type string
	var vm_index int
	// var vm_storage
	// var vm_storage string
	var vm_misc string
	var stdout []byte

	for index, vm := range vm_list.vmName {
		vm_index = index + 1
		vm_status = vmStatusCheck(vm).vmStatusIcons
		vm_dataset = vm_list.vmDataset[index]
		vm_resources = "CPUs: " + strconv.Itoa(VmConfig(vm).Cpus) + "\nRAM: " + strconv.Itoa(VmConfig(vm).Ram) + "G"
		vm_vnc = "Port: " + strconv.Itoa(VmConfig(vm).VncPort) + "\nPwd: " + VmConfig(vm).VncPassword
		vm_networks = VmConfig(vm).Networks[0].InterfaceName + ": " + VmConfig(vm).Networks[0].InterfaceIpAddress
		vm_misc = "OS: " + vm_os_type + "\nUptime: 00:00" + "\nParent: " + VmConfig(vm).ParentHost

		//Storage
		vm_storage_full_size, _ := os.Stat(VmConfig(vm).Storage[0].DiskLocation)
		vm_storage_provisioned := vm_storage_full_size.Size()
		command_ := "du /zroot/vm-encrypted/test-vm-1/disk0.img | awk '{ print $1 }'"
		command := exec.Command("bash", "-c", command_)
		stdout, _ = command.Output()

		// OS Types hot replacement
		vm_os_type = strings.ReplaceAll(VmConfig(vm).OsType, "debian11", "Debian 11")

		outputTable.AppendRow([]interface{}{
			vm_index,
			vm,
			vm_status,
			vm_dataset,
			vm_resources,
			vm_vnc,
			vm_networks,
			vm_storage_provisioned,
			vm_misc})
		outputTable.AppendSeparator()
	}

	var total_number_of_vms = strconv.Itoa(len(vm_list.vmName))
	outputTable.AppendFooter(table.Row{"", "total vms: " + total_number_of_vms})

	fmt.Println(stdout)
	outputTable.SetStyle(table.StyleLight)
	outputTable.Render()
}

type vmListStruct struct {
	vmName    []string
	vmDataset []string
}

func vmList(plain ...bool) vmListStruct {
	var datasetsList_var = datasetsList()
	var folder_to_scan string
	var vm_list = vmListStruct{}

	//Form VM list from all available datasets
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
			}

			var _, new_config_file_exists_error = os.Stat(vm_folder_full_path + "/vm.conf")
			if new_config_file_exists_error == nil {
				vm_list.vmName = append(vm_list.vmName, vm_folder_name)
			}
		}
	}

	//Sort the VM list
	natsort.Sort(vm_list.vmName)

	//Form the list of dataset names
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

type vmStatusCheckStruct struct {
	vmLive        bool
	vmEncrypted   bool
	vmStatusIcons string
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

type datasetsListStruct struct {
	Datasets []struct {
		Name       string `yaml:"name"`
		Mount_path string `yaml:"mount_path"`
		Zfs_path   string `yaml:"zfs_path"`
		Encrypted  bool   `yaml:"encrypted"`
		Type       string `yaml:"type"`
	}
}

func datasetsList() datasetsListStruct {
	var conf_datasets_file, conf_datasets_error = os.ReadFile("conf_datasets.yaml")

	if conf_datasets_error != nil {
		panic(conf_datasets_error)
	}

	var datasetsList_var datasetsListStruct

	err := yaml.Unmarshal([]byte(conf_datasets_file), &datasetsList_var)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return datasetsList_var
}
