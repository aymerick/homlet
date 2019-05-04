package homlet

import (
	"io"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
)

const (
	BAUD    = 57600
	LF_CHAR = 10
)

// Hardware is the serial connection with homlet device master
type Hardware struct {
	ser io.ReadWriteCloser
}

// Open instanciates and initializes hardware
func Open(path string) (*Hardware, error) {
	config := &serial.Config{Name: path, Baud: BAUD}
	ser, err := serial.OpenPort(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open serial port")
	}

	result := Hardware{ser: ser}
	return &result, nil
}

// Close stops hardware
func (hdw *Hardware) Close() {
	hdw.ser.Close()
}

// Read returns a packet
func (hdw *Hardware) Read() (*Packet, error) {
	for {
		line := []byte{}
		lastRead := make([]byte, 1)

		// read byte by byte until the Line Feed character
		// FIXME optimize that
		for lastRead[0] != LF_CHAR {
			n, err := hdw.ser.Read(lastRead)
			if err != nil {
				return nil, err
			}

			if n != 1 {
				return nil, errors.New("no data read")
			}

			line = append(line, lastRead[0])
		}

		packet, err := Parse(string(line))
		if err == errLineFormat {
			continue
		}

		return packet, err
	}
}

// ReadPackets read packets and send them to returned channel
func (hdw *Hardware) ReadPackets() chan *Packet {
	result := make(chan *Packet)

	go func(hdw *Hardware, c chan *Packet) {
		for {
			packet, err := hdw.Read()
			if err != nil {
				log.Errorf("Failed to read data: %s", err)
				close(c)
				return
			}

			c <- packet
		}
	}(hdw, result)

	return result
}
