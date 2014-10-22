package serial

import (
	"bufio"
	"io"
	"log"
	"sync"

	"github.com/aymerick/homlet"
	"github.com/tarm/goserial"
)

type SerialHardware struct {
	homlet.Hardware

	port string
	baud int

	serialPort io.ReadWriteCloser
}

func NewSerialHardware(kind string, name string, port string, baud int) *SerialHardware {
	return &SerialHardware{
		Hardware: *homlet.NewHardware(kind, name),
		port:     port,
		baud:     baud,
	}
}

func (self *SerialHardware) Port() string {
	return self.port
}

func (self *SerialHardware) Baud() int {
	return self.baud
}

// Prints debug message
func (self *SerialHardware) Debug() {
	log.Printf("[%v] %v on %v at %v bauds", self.Name(), self.Kind(), self.Port(), self.Baud())
}

// Starts hardware
func (self *SerialHardware) Start(wg *sync.WaitGroup) {
	log.Printf("[%v] %v > Starting (TODO)", self.Name(), self.Kind())

	var err error

	self.serialPort, err = serial.OpenPort(&serial.Config{Name: self.Port(), Baud: self.Baud()})
	if err != nil {
		log.Fatal(err)
	}

	// NOTE: serial read is blocking, so no way to do a gracefull shutdown
	// wg.Add(1)

	go func() {
		// defer wg.Done()

		scanner := bufio.NewScanner(self.serialPort)
		for scanner.Scan() {
			txt := scanner.Text()
			log.Printf(txt)
			// @todo handle data !
		}

		if err := scanner.Err(); err != nil {
			log.Printf("Failed to read serial port: %v", err)
		}
	}()
}

// Stops hardware
func (self *SerialHardware) Stop() {
	// NOP ... serial read is blocking, so no way to do a gracefull shutdown
}
