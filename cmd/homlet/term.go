package main

import (
	"io/ioutil"

	"github.com/aymerick/homlet"
	"github.com/aymerick/homlet/pkg/term"
	"github.com/gizak/termui/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
	hdw, err := homlet.Open("/dev/tty.usbserial-A1014IM4")
	if err != nil {
		log.Fatalf("Failed to open serial connection: %v", err)
	}
	defer hdw.Close()

	// do not mess UI with log messages
	log.SetOutput(ioutil.Discard)

	// run UI
	term.NewUI(readPackets(hdw)).Run()
}

func readPackets(hdw *homlet.Hardware) chan *homlet.Packet {
	result := make(chan *homlet.Packet)

	go func(hdw *homlet.Hardware, c chan *homlet.Packet) {
		for {
			packet, err := hdw.Read()
			if err != nil {
				log.Errorf("Failed to read data: %s", err)
				close(c)
				return
			}

			c <- packet
		}
	}(hdw, result)

	return result
}
