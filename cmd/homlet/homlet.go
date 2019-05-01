package main

import (
	"github.com/aymerick/homlet"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)

	hdw, err := homlet.Open("/dev/tty.usbserial-A1014IM4")
	if err != nil {
		log.Error("Failed to open serial connection")
		return
	}

	defer hdw.Close()

	for {
		packet, err := hdw.Read()
		if err != nil {
			log.Errorf("Failed to read data: %s", err)
			return
		}

		// TODO handle packet
		log.Infof("Received: %s", packet)
	}
}
