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

func (hardware *SerialHardware) Port() string {
	return hardware.port
}

func (hardware *SerialHardware) Baud() int {
	return hardware.baud
}

func (hardware *SerialHardware) SetProcessor(processor SerialProcesserInterface) {
	hardware.processor = processor
}

/**
 * HardwareInterface
 */

// Implements HardwareInterface, overwrites Hardware#Debug
func (hardware *SerialHardware) Debug() {
	log.Printf("[%v] %v on %v at %v bauds", hardware.Name(), hardware.Kind(), hardware.Port(), hardware.Baud())
}

// Implements HardwareInterface
func (hardware *SerialHardware) Start(wg *sync.WaitGroup) {
	log.Printf("[%v] %v > Starting (TODO)", hardware.Name(), hardware.Kind())

	var err error

	hardware.serialPort, err = serial.OpenPort(&serial.Config{Name: hardware.Port(), Baud: hardware.Baud()})
	if err != nil {
		log.Fatal(err)
	}

	// NOTE: serial read is blocking, so no way to do a gracefull shutdown
	// wg.Add(1)

	go func() {
		// NOTE: serial read is blocking, so no way to do a gracefull shutdown
		// defer wg.Done()

		scanner := bufio.NewScanner(hardware.serialPort)
		for scanner.Scan() {
			txt := scanner.Text()

			if hardware.processor == nil {
				log.Fatal(errors.New("No processor set to handle incoming data"))
			}

			hardware.processor.ProcessLine(txt)
		}

		if err := scanner.Err(); err != nil {
			log.Printf("Failed to read serial port: %v", err)
		}
	}()
}

// Implements HardwareInterface
func (hardware *SerialHardware) Stop() {
	// NOOP ... serial read is blocking, so no way to do a gracefull shutdown
}
