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

func (self *Homlet) Hardwares() *hardwares {
	return self.hardwares
}

// Adds hardware to app
func (self *Homlet) AddHardware(hardware HardwareInterface) {
	*self.hardwares = append(*self.Hardwares(), hardware)
}

// Start the engine
func (self *Homlet) Start() {
	self.Hardwares().Start(self.hardwaresWG)

	self.dispatcher.start()
}

// Stop the engine
func (self *Homlet) Stop() {
	self.dispatcher.stop()

	self.Hardwares().Stop(self.hardwaresWG)
}
