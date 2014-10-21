package rfxcom

import (
	"log"

	"github.com/aymerick/homlet"
)

type RfxcomHardware struct {
	homlet.Hardware
}

func NewRfxcomHardware(name string) *RfxcomHardware {
	return &RfxcomHardware{
		Hardware: *homlet.NewHardware("rfxcom", name),
	}
}

// Starts hardware
func (self *RfxcomHardware) Start() {
	log.Printf("[%v] %v > Starting (TODO)", self.Kind(), self.Name())
}

// Stops hardware
func (self *RfxcomHardware) Stop() {
	log.Printf("[%v] %v > Stopping (TODO)", self.Kind(), self.Name())
}
