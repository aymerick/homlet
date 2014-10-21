package main

import (
	"log"

	"github.com/aymerick/homlet"
	"github.com/aymerick/homlet/hardware/homlet_node"
	"github.com/aymerick/homlet/hardware/rfxcom"
)

func main() {
	log.Printf("Cooking Homletd")

	// Init application
	app := homlet.NewHomlet()

	// @todo FINISH THAT !
	app.AddHardware(homlet_node.NewHomletNodeHardware("Jeelink"))
	app.AddHardware(rfxcom.NewRfxcomHardware("RFXtrx433E"))

	// Start daemon
	daemon := NewHomletd(app)

	daemon.run()

	log.Printf("Homletd ended")
}
