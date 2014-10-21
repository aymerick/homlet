package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aymerick/homlet"
)

// Daemon
type Homletd struct {
	// application
	app *homlet.Homlet

	// Are we shutting down ?
	shuttingDown bool

	// Channel that receives signals
	sigChan chan os.Signal
}

// Invoc a new daemon
func NewHomletd(app *homlet.Homlet) *Homletd {
	return &Homletd{
		app:          app,
		shuttingDown: false,
		sigChan:      make(chan os.Signal, 1),
	}
}

func (self *Homletd) isShuttingDown() bool {
	return self.shuttingDown
}

func (self *Homletd) run() {
	// wait for sigint
	signal.Notify(self.sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGUSR1)

	for !self.isShuttingDown() {
		select {
		case sig := <-self.sigChan:
			switch sig {
			case syscall.SIGHUP:
				log.Println("SIGHUP - Reload initiated (TODO)")
			case syscall.SIGINT, syscall.SIGTERM:
				log.Println("SIGINT or SIGTERM received - Shutdown initiated")
				self.stop()
			case syscall.SIGUSR1:
				log.Println("SIGUSR1 received - TODO")
			}
		}
	}
}

func (self *Homletd) stop() {
	self.shuttingDown = true
}
