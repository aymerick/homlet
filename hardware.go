package homlet

import (
	"log"
	"sync"
)

type HardwareInterface interface {
	Kind() string
	Name() string

	Debug()

	Start(wg *sync.WaitGroup)
	Stop()
}

// Base hardware
type Hardware struct {
	kind string
	name string

	wg *sync.WaitGroup
}

type hardwares []HardwareInterface

func NewHardware(kind string, name string) *Hardware {
	return &Hardware{
		kind: kind,
		name: name,
	}
}

// Implements HardwareInterface
func (self *Hardware) Kind() string {
	return self.kind
}

// Implements HardwareInterface
func (self *Hardware) Name() string {
	return self.name
}

// Implements HardwareInterface
func (self *Hardware) Debug() {
	log.Printf("[%v] %v", self.Kind(), self.Name())
}

/**
 * hardwares
 */

// Starts all hardwares
func (col *hardwares) Start(wg *sync.WaitGroup) {
	for _, hardware := range *col {
		hardware.Start(wg)
	}
}

// Stops all hardwares
func (col *hardwares) Stop(wg *sync.WaitGroup) {
	for _, hardware := range *col {
		hardware.Stop()
	}

	// wait for all hardwares to stop
	wg.Wait()
}

// Debugs all hardwares
func (col *hardwares) Debug() {
	for _, hardware := range *col {
		hardware.Debug()
	}
}
