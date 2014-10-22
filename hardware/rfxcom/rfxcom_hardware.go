package rfxcom

import "github.com/aymerick/homlet/hardware/serial"

type RfxcomHardware struct {
	serial.SerialHardware
}

func NewRfxcomHardware(name string, port string) *RfxcomHardware {
	return &RfxcomHardware{
		SerialHardware: *serial.NewSerialHardware("rfxcom", name, port, 4800),
	}
}
