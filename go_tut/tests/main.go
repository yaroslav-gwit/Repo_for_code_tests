package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type datasetZfs struct {
	Datasets []struct {
		Name       string `yaml:"name"`
		Mount_path string `yaml:"mount_path"`
		Zfs_path   string `yaml:"zfs_path"`
		Encrypted  bool   `yaml:"encrypted"`
	}
}

func main() {
	var conf_datasets_file, conf_datasets_error = os.ReadFile("conf_datasets.yaml")

	if conf_datasets_error != nil {
		panic(conf_datasets_error)
	}

	var datasetZfs_var = datasetZfs{}

	err := yaml.Unmarshal([]byte(conf_datasets_file), &datasetZfs_var)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	for _, dataset := range datasetZfs_var.Datasets {
		fmt.Println(dataset.Mount_path)
	}
}
