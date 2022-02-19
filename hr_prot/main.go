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
	DiskSize        int
	DiskFreeGig     int
	DiskFreePercent int
}

func main() {
	VmUptime("vm")
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
	var vm_misc string

	for index, vmname := range vm_list {
		// OS Types hot replacement
		vm_os_type = VmConfig(vmname).OsType
		vm_os_type = strings.ReplaceAll(vm_os_type, "debian11", "Debian 11")
		vm_os_type = strings.ReplaceAll(vm_os_type, "freebsd13ufs", "FreeBSD 13 UFS")

		vm_index = index + 1
		vm_status = vmStatusCheck(vmname).vmStatusIcons
		vm_dataset = VmDatasetCheck(vmname).Name
		vm_resources = "CPUs: " + strconv.Itoa(VmConfig(vmname).Cpus) + "\nRAM: " + strconv.Itoa(VmConfig(vmname).Ram) + "G"
		vm_vnc = "Port: " + strconv.Itoa(VmConfig(vmname).VncPort) + "\nPwd: " + VmConfig(vmname).VncPassword
		vm_networks = VmConfig(vmname).Networks[0].InterfaceName + ": " + VmConfig(vmname).Networks[0].InterfaceIpAddress
		vm_misc = "OS: " + vm_os_type + "\nUptime: 00:00" + "\nParent: " + VmConfig(vmname).ParentHost

		//Storage
		vm_disk_location := VmConfig(vmname).Storage[0].DiskFolder + VmConfig(vmname).Storage[0].DiskImage
		vm_storage_full_size, _ := os.Stat(vm_disk_location)
		vm_storage_provisioned := vm_storage_full_size.Size()
		vm_storage_provisioned = vm_storage_provisioned / 1024 / 1024 / 1024
		cmd := "du " + vm_disk_location + " | awk '{ print $1 }'"
		var out, _ = exec.Command("bash", "-c", cmd).Output()
		vm_storage_used_ := strings.ReplaceAll(string(out), "\n", "")
		vm_storage_used, _ := strconv.Atoi(vm_storage_used_)
		vm_storage_used = vm_storage_used / 1024 / 1024
		vm_storage := VmConfig(vmname).Storage[0].DiskName + ": " + strconv.Itoa(vm_storage_used) + "G/" + strconv.Itoa(int(vm_storage_provisioned)) + "G"

		outputTable.AppendRow([]interface{}{
			vm_index,
			vmname,
			vm_status,
			vm_dataset,
			vm_resources,
			vm_vnc,
			vm_networks,
			vm_storage,
			vm_misc})
		outputTable.AppendSeparator()
	}

	var total_number_of_vms = strconv.Itoa(len(vm_list))
	outputTable.AppendFooter(table.Row{"", "total vms: " + total_number_of_vms})
	outputTable.SetStyle(table.StyleLight)
	outputTable.Render()
}

func vmList(plain ...bool) []string {
	var datasetsList_var = datasetsList()
	var folder_to_scan string
	// var vm_list = vmListStruct{}
	var vm_list []string

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

			var _, file_exists_error = os.Stat(vm_folder_full_path + "/conf_vm.yaml")
			if file_exists_error == nil {
				vm_list = append(vm_list, vm_folder_name)
			}
		}
	}

	//Sort the VM list
	natsort.Sort(vm_list)

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
	var dataset = VmDatasetCheck(vmname)

	if dataset.Encrypted {
		vmStatusIcons = vmStatusIcons + vm_is_encrypted
		vmStatusCheckStruct_var.vmEncrypted = true
	} else {
		vmStatusCheckStruct_var.vmEncrypted = false
	}

	vmStatusCheckStruct_var.vmStatusIcons = vmStatusIcons
	return vmStatusCheckStruct_var
}

//Dataset LIST section
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

//SINGLE Dataset section

type datasetStruct struct {
	Name       string `yaml:"name"`
	Mount_path string `yaml:"mount_path"`
	Zfs_path   string `yaml:"zfs_path"`
	Encrypted  bool   `yaml:"encrypted"`
	Type       string `yaml:"type"`
}

func VmDatasetCheck(vmname string) datasetStruct {
	//Load datasets list file
	var conf_datasets_file, conf_datasets_error = os.ReadFile("conf_datasets.yaml")
	if conf_datasets_error != nil {
		panic(conf_datasets_error)
	}

	//Init and populate the variable for list of datasets
	var vm_dataset_list datasetsListStruct
	err := yaml.Unmarshal([]byte(conf_datasets_file), &vm_dataset_list)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	//Iterate over the list of datasets and check where VM exists
	var vm_dataset datasetStruct
	for _, dataset := range vm_dataset_list.Datasets {
		folder_to_scan := dataset.Mount_path
		_, vm_in_dataset_error := os.Stat(folder_to_scan + vmname)
		if vm_in_dataset_error == nil {
			vm_dataset.Name = dataset.Name
			vm_dataset.Mount_path = dataset.Mount_path
			vm_dataset.Zfs_path = dataset.Zfs_path
			vm_dataset.Encrypted = dataset.Encrypted
			vm_dataset.Type = dataset.Type
		}
	}

	return vm_dataset
}

//VM Config Section
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
		DiskName      string `yaml:"disk_name"`
		DiskType      string `yaml:"disk_type"`
		DiskDriveType string `yaml:"disk_drive_type"`
		DiskFolder    string `yaml:"disk_folder"`
		DiskImage     string `yaml:"disk_image"`
	}
}

func VmConfig(vmname string) vmConfigStruct {
	var conf_vm_file []byte
	var vm_dataset = VmDatasetCheck(vmname)

	var _conf_vm_file, conf_vm_error = os.ReadFile(vm_dataset.Mount_path + vmname + "/conf_vm.yaml")
	if conf_vm_error == nil {
		conf_vm_file = _conf_vm_file
	} else {
		panic("Can't find config file!")
	}

	var vmConfigStruct_var vmConfigStruct

	err := yaml.Unmarshal([]byte(conf_vm_file), &vmConfigStruct_var)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return vmConfigStruct_var
}

func VmUptime(vmname string) string {
	cmd := "ps axwww -o etime,command > /tmp/bhyve_vms_uptime.txt"
	var _, _ = exec.Command("bash", "-c", cmd).Output()
	// vm_storage_used_ := strings.ReplaceAll(string(out), "\n", "")

	var vm_uptime_file, vm_uptime_error = os.ReadFile("/tmp/bhyve_vms_uptime.txt")
	if vm_uptime_error != nil {
		panic("Can't find config file!")
	}

	// vm_storage_used_ := strings.Split(string(out), "\n")
	blah := strings.Split(string(vm_uptime_file), "\n")
	fmt.Println(blah)
	return "asdf"
}
