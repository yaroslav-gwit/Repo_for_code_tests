package main

import (
	"fmt"
	"reflect"

	"github.com/spf13/viper"
)

type BlaBla []interface {
	datasetsViperInterface()
}

type datasetsListStruct struct {
	Datasets []struct {
		Name       string
		Mount_path string
		Zfs_path   string
		Encrypted  bool
		Type       string
	}
}

func main() {
	datasetsViper()
}

func datasetsViper() interface{} {
	viper.SetConfigName("conf_datasets")
	// viper.AddConfigPath("/etc/appname/")
	// viper.AddConfigPath("$HOME/.appname")
	viper.AddConfigPath(".")
	viper_err := viper.ReadInConfig()
	if viper_err != nil {
		panic(fmt.Errorf("fatal error config file: %w", viper_err))
	}

	datasetsViper := viper.Get("datasets")

	fmt.Println(reflect.TypeOf(datasetsViper))
	fmt.Println(datasetsViper)

	viper.WriteConfigAs("./viper_config.yaml")

	return datasetsViper
}
