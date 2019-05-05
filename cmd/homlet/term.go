package main

import (
	"io/ioutil"

	"github.com/aymerick/homlet"
	"github.com/aymerick/homlet/pkg/term"
	"github.com/gizak/termui/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var termCmd = &cobra.Command{
	Use:   "term",
	Short: "Launch the terminal UI",
	Run:   runTerm,
}

func init() {
	rootCmd.AddCommand(termCmd)
}

func runTerm(cmd *cobra.Command, args []string) {
	// init UI
	if err := termui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer termui.Close()

	// init hardware
	hdw, err := homlet.Open(viper.GetString("serial"))
	if err != nil {
		log.Fatalf("Failed to open serial connection: %v", err)
	}
	defer hdw.Close()

	// do not mess UI with log messages
	log.SetOutput(ioutil.Discard)

	// get devices settings
	settings, err := devicesSettings()
	if err != nil {
		log.Fatalf("Failed to fetch devices settings: %v", err)
	}

	// read packets
	packets := hdw.ReadPackets()

	// run UI
	term.NewUI(packets, settings).Run()
}
