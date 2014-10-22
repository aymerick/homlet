package homlet_node

import (
	"log"

	"github.com/aymerick/homlet/hardware/serial"
)

type HomletNodeHardware struct {
	serial.SerialHardware
}

func NewHomletNodeHardware(name string, port string) *HomletNodeHardware {
	result := &HomletNodeHardware{
		SerialHardware: *serial.NewSerialHardware("homlet_node", name, port, 57600),
	}

	result.SerialHardware.SetProcessor(result)

	return result
}

// Implements SerialProcesserInterface
func (self *HomletNodeHardware) ProcessLine(line string) {
	log.Printf("homlet_node > %v", line)
}
