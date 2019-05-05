package main

import (
	"fmt"
	"os"

	"github.com/aymerick/homlet"
	"github.com/spf13/viper"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func devicesSettings() ([]*homlet.DeviceSettings, error) {
	result := []*homlet.DeviceSettings{}
	if err := viper.UnmarshalKey("devices", &result); err != nil {
		return nil, err
	}
	return result, nil
}
