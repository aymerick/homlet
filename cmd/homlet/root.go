package main

import (
	"fmt"
	"os"
	"runtime"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "homlet",
	Short: "homlet is a DIY domotic system",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// set log level
		l := viper.GetString("loglevel")
		level, err := log.ParseLevel(l)
		if err != nil {
			panic(fmt.Sprintf("unexpected log level '%s': %s", l, err))
		}
		log.SetLevel(level)
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default \"/etc/homlet/config.toml\" or \"$HOME/.homlet/config.toml)\"")

	rootCmd.PersistentFlags().StringP("loglevel", "l", "info", "Log level")
	if err := viper.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel")); err != nil {
		panic(err)
	}

	rootCmd.PersistentFlags().String("serial", defaultSerial(), "Serial port")
	if err := viper.BindPFlag("serial", rootCmd.PersistentFlags().Lookup("serial")); err != nil {
		panic(err)
	}
}

func defaultSerial() string {
	switch runtime.GOOS {
	case "darwin":
		return "/dev/tty.usbserial-A1014IM4"
	case "linux":
		switch runtime.GOARCH {
		case "arm", "arm64":
			// Raspberry on FTDI serial (cf. http://jeelabs.org/2012/09/20/serial-hookup-jeenode-to-raspberry-pi/)
			return "/dev/ttyAMA0"
		default:
			return "/dev/ttyUSB0"
		}
	default:
		panic(fmt.Sprintf("unsupported OS: %s", runtime.GOOS))
	}
}

func initConfig() {
	// setup environment variables
	viper.SetEnvPrefix("homlet")
	viper.AutomaticEnv()

	// setup config file
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")         // name of config file (without extension)
		viper.AddConfigPath("/etc/homlet/")   // path to look for the config file in
		viper.AddConfigPath("$HOME/.homlet/") // call multiple times to add many search paths
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Println("Can't read config:", err)
			os.Exit(1)
		}
	}
}
