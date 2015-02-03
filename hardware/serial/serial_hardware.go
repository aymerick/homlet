package serial

import (
	"bufio"
	"errors"
	"io"
	"log"
	"sync"

	"github.com/aymerick/homlet"
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

func (serial *SerialHardware) Port() string {
	return serial.port
}

func (serial *SerialHardware) Baud() int {
	return serial.baud
}

func (serial *SerialHardware) SetProcessor(processor SerialProcesserInterface) {
	serial.processor = processor
}

/**
 * HardwareInterface
 */

// Implements HardwareInterface, overwrites Hardware#Debug
func (serial *SerialHardware) Debug() {
	log.Printf("[%v] %v on %v at %v bauds", serial.Name(), serial.Kind(), serial.Port(), serial.Baud())
}

// Implements HardwareInterface
func (serial *SerialHardware) Start(wg *sync.WaitGroup) {
	log.Printf("[%v] %v > Starting (TODO)", serial.Name(), serial.Kind())

	var err error

	serial.serialPort, err = serial.OpenPort(&serial.Config{Name: serial.Port(), Baud: serial.Baud()})
	if err != nil {
		log.Fatal(err)
	}

	// NOTE: serial read is blocking, so no way to do a gracefull shutdown
	// wg.Add(1)

	go func() {
		// NOTE: serial read is blocking, so no way to do a gracefull shutdown
		// defer wg.Done()

		scanner := bufio.NewScanner(serial.serialPort)
		for scanner.Scan() {
			txt := scanner.Text()

			if serial.processor == nil {
				log.Fatal(errors.New("No processor set to handle incoming data"))
			}

			serial.processor.ProcessLine(txt)
		}

		if err := scanner.Err(); err != nil {
			log.Printf("Failed to read serial port: %v", err)
		}
	}()
}

// Implements HardwareInterface
func (serial *SerialHardware) Stop() {
	// NOOP ... serial read is blocking, so no way to do a gracefull shutdown
}
