package rfxcom

import "github.com/aymerick/homlet"

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
	// @todo
	panic("not implemented")
}

// Stops hardware
func (self *RfxcomHardware) Stop() {
	// @todo
	panic("not implemented")
}
