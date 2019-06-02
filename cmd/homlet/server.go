package main

import (
	"github.com/aymerick/homlet"
	"github.com/aymerick/homlet/pkg/domoticz"
	"github.com/aymerick/homlet/pkg/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Launch the server",
	Run:   runServer,
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func runServer(cmd *cobra.Command, args []string) {
	// init hardware
	hdw, err := homlet.Open(viper.GetString("serial"))
	if err != nil {
		log.Fatalf("Failed to open serial connection: %v", err)
	}
	defer hdw.Close()

	// get devices settings
	settings, err := devicesSettings()
	if err != nil {
		log.Fatalf("Failed to fetch devices settings: %v", err)
	}

	// read packets
	packets := hdw.ReadPackets()

	// init server
	server := server.New(packets, settings)

	if url := viper.GetString("domoticz.url"); url != "" {
		handler := domoticz.New(viper.GetInt("domoticz.hardwareId"), url)
		server.SetDomoticz(handler)
	}

	// run server
	server.Run()
}
