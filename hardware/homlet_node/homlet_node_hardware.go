package homlet_node

import "github.com/aymerick/homlet/hardware/serial"

type HomletNodeHardware struct {
	serial.SerialHardware
}

func NewHomletNodeHardware(name string, port string) *HomletNodeHardware {
	return &HomletNodeHardware{
		SerialHardware: *serial.NewSerialHardware("homlet_node", name, port, 57600),
	}
}
