package homlet_node

import (
	"log"

	"github.com/aymerick/homlet"
)

type HomletNodeHardware struct {
	homlet.Hardware
}

func NewHomletNodeHardware(name string) *HomletNodeHardware {
	return &HomletNodeHardware{
		Hardware: *homlet.NewHardware("homlet_node", name),
	}
}

// Starts hardware
func (self *HomletNodeHardware) Start() {
	log.Printf("[%v] %v > Starting (TODO)", self.Kind(), self.Name())
}

// Stops hardware
func (self *HomletNodeHardware) Stop() {
	log.Printf("[%v] %v > Stopping (TODO)", self.Kind(), self.Name())
}
