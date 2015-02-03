package homlet

import "sync"

// Master object encapsulating the Homlet application conf and states
type Homlet struct {
	// All running hardwares
	hardwares *hardwares

	// Wait group when stopping hardwares
	hardwaresWG *sync.WaitGroup

	// Messages dispatcher
	dispatcher *Dispatcher
}

// Creates a new Homlet app
func NewHomlet() *Homlet {
	return &Homlet{
		hardwares:   &hardwares{},
		hardwaresWG: &sync.WaitGroup{},
		dispatcher:  NewDispatcher(),
	}
}

func (homlet *Homlet) Hardwares() *hardwares {
	return homlet.hardwares
}

// Adds hardware to app
func (homlet *Homlet) AddHardware(hardware HardwareInterface) {
	*homlet.hardwares = append(*homlet.Hardwares(), hardware)
}

// Start the engine
func (homlet *Homlet) Start() {
	homlet.Hardwares().Start(homlet.hardwaresWG)

	homlet.dispatcher.start()
}

// Stop the engine
func (homlet *Homlet) Stop() {
	homlet.dispatcher.stop()

	homlet.Hardwares().Stop(homlet.hardwaresWG)
}
