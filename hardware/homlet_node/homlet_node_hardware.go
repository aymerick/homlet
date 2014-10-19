package homlet_node

import "github.com/aymerick/homlet"

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
	// @todo
	panic("not implemented")
}

// Stops hardware
func (self *HomletNodeHardware) Stop() {
	// @todo
	panic("not implemented")
}
