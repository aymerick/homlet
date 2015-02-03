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

func (homletd *Homletd) isShuttingDown() bool {
	return homletd.shuttingDown
}

func (homletd *Homletd) run() {
	// start app
	homletd.app.Start()

	// wait for signals
	signal.Notify(homletd.sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGUSR1)

	for !homletd.isShuttingDown() {
		select {
		case sig := <-homletd.sigChan:
			switch sig {
			case syscall.SIGHUP:
				log.Println("SIGHUP - Reload initiated (TODO)")
			case syscall.SIGINT, syscall.SIGTERM:
				log.Println("SIGINT or SIGTERM received - Shutdown initiated")
				homletd.stop()
			case syscall.SIGUSR1:
				log.Println("SIGUSR1 received - Printing status (TODO)")
				// Debug hardwares for now
				homletd.app.Hardwares().Debug()
			}
		}
	}

	// stop app
	homletd.app.Stop()

	log.Println("Shutdown complete")
}

func (homletd *Homletd) stop() {
	homletd.shuttingDown = true
}
