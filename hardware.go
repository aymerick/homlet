package homlet

import (
	"io"

	"github.com/pkg/errors"
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

func Open(path string) (*Hardware, error) {
	config := &serial.Config{Name: path, Baud: BAUD}
	ser, err := serial.OpenPort(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open serial port")
	}

	result := Hardware{ser: ser}
	return &result, nil
}

func (hdw *Hardware) Close() {
	hdw.ser.Close()
}

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
