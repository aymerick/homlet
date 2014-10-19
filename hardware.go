package homlet

import "log"

type HardwareInterface interface {
	Kind() string
	Name() string

	Debug()

	Start()
	Stop()
}

// Base hardware
type Hardware struct {
	kind string
	name string
}

type hardwares []HardwareInterface

/**
 * Hardware
 */

func NewHardware(kind string, name string) *Hardware {
	return &Hardware{
		kind: kind,
		name: name,
	}
}

// Get hardware kind
func (self *Hardware) Kind() string {
	return self.kind
}

// Get hardware name
func (self *Hardware) Name() string {
	return self.name
}

// Prints debug message
func (self *Hardware) Debug() {
	log.Printf("[%v] %v", self.Kind(), self.Name())
}

/**
 * hardwares
 */

// Starts all hardwares
func (col *hardwares) Start() {
	for _, hardware := range *col {
		hardware.Start()
	}
}

// Stops all hardwares
func (col *hardwares) Stop() {
	for _, hardware := range *col {
		hardware.Stop()
	}
}

// Debugs all hardwares
func (col *hardwares) Debug() {
	for _, hardware := range *col {
		hardware.Debug()
	}
}
