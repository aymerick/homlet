package main

import (
	"log"

	"github.com/aymerick/homlet"
	"github.com/aymerick/homlet/hardware/homlet_node"
)

func main() {
	log.Printf("Cooking Homletd")

	// Init application
	app := homlet.NewHomlet()

	// @todo FINISH THAT !
	app.AddHardware(homlet_node.NewHomletNodeHardware("Jeelink", "/dev/tty.usbserial-A1014IM4"))
	// app.AddHardware(rfxcom.NewRfxcomHardware("RFXtrx433E", "TODO"))

	// debug
	app.Hardwares().Debug()

	// Start daemon
	daemon := NewHomletd(app)

	daemon.run()

	log.Printf("Homletd ended")
}
