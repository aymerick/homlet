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
func (hardware *Hardware) Kind() string {
	return hardware.kind
}

// Implements HardwareInterface
func (hardware *Hardware) Name() string {
	return hardware.name
}

// Implements HardwareInterface
func (hardware *Hardware) Debug() {
	log.Printf("[%v] %v", hardware.Kind(), hardware.Name())
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
