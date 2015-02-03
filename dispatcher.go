package homlet

import "log"

// Dispatcher
type Dispatcher struct {
	// Channel to stop dispatcher
	stopChan chan bool

	// Channel that receives status messages from hardwares
	statusChan chan *StatusMessage

	// Channel that receives command messages for hardwares
	cmdChan chan *CommandMessage
}

// Creates a new dispatcher
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		stopChan:   make(chan bool),
		statusChan: make(chan *StatusMessage, 1),
		cmdChan:    make(chan *CommandMessage, 1),
	}
}

// Start Dispatcher
func (dispatcher *Dispatcher) start() {
	go func() {
		stop := false

		for !stop {
			select {
			case <-dispatcher.stopChan:
				log.Printf("[Dispatcher] Stop received")
				stop = true

			case statusMsg := <-dispatcher.statusChan:
				log.Printf("[Dispatcher] Status message received (TODO): %v", statusMsg)

			case cmdMsg := <-dispatcher.cmdChan:
				log.Printf("[Dispatcher] Command message received (TODO): %v", cmdMsg)
			}
		}

		log.Printf("[Dispatcher] Stopped")

		// ok we are done
		close(dispatcher.stopChan)
	}()

	log.Printf("[Dispatcher] Started")
}

// Stop Dispatcher
func (dispatcher *Dispatcher) stop() {
	// stop go routine
	dispatcher.stopChan <- true

	// wait for go routine to end
	<-dispatcher.stopChan
}
