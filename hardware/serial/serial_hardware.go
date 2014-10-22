package serial

import (
	"bufio"
	"errors"
	"io"
	"log"
	"sync"

	"github.com/aymerick/homlet"
	"github.com/tarm/goserial"
)

type SerialProcesserInterface interface {
	ProcessLine(string)
}

type SerialHardware struct {
	homlet.Hardware

	port string
	baud int

	serialPort io.ReadWriteCloser
	processor  SerialProcesserInterface
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

func (self *SerialHardware) SetProcessor(processor SerialProcesserInterface) {
	self.processor = processor
}

/**
 * HardwareInterface
 */

// Implements HardwareInterface, overwrites Hardware#Debug
func (self *SerialHardware) Debug() {
	log.Printf("[%v] %v on %v at %v bauds", self.Name(), self.Kind(), self.Port(), self.Baud())
}

// Implements HardwareInterface
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

			if self.processor == nil {
				log.Fatal(errors.New("No processor set to handle incoming data"))
			}

			self.processor.ProcessLine(txt)
		}

		if err := scanner.Err(); err != nil {
			log.Printf("Failed to read serial port: %v", err)
		}
	}()
}

// Implements HardwareInterface
func (self *SerialHardware) Stop() {
	// NOP ... serial read is blocking, so no way to do a gracefull shutdown
}
