package homlet

// Master object encapsulating the Homlet application conf and states
type Homlet struct {
	// All running hardwares
	hardwares *hardwares
}

// Creates a new Homlet app
func NewHomlet() *Homlet {
	return &Homlet{
		hardwares: &hardwares{},
	}
}

func (self *Homlet) Hardwares() *hardwares {
	return self.hardwares
}

// Adds hardware to app
func (self *Homlet) AddHardware(hardware HardwareInterface) {
	// @todo Add a sync.Mutex the day we need it
	*self.hardwares = append(*self.Hardwares(), hardware)
}

// Start the engine
func (self *Homlet) Start() {
	self.Hardwares().Start()
}

// Stop the engine
func (self *Homlet) Stop() {
	self.Hardwares().Stop()
}
