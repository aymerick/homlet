package main

import (
	"log"

	"github.com/aymerick/homlet"
	"github.com/aymerick/homlet/hardware/homlet_node"
	"github.com/aymerick/homlet/hardware/rfxcom"
)

func main() {
	log.Printf("Cooking Homlet")

	app := homlet.NewHomlet()

	// @todo FINISH THAT !
	app.AddHardware(homlet_node.NewHomletNodeHardware("Jeelink"))
	app.AddHardware(rfxcom.NewRfxcomHardware("RFXtrx433E"))

	// Debug hardwares
	app.Hardwares().Debug()

	log.Printf("Homlet ended")
}
