package homlet

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var errLineFormat = errors.New("erroneous line received")

// Packet received from a device
type Packet struct {
	At          time.Time
	Corrections map[Sensor]float64
	Device      Device
	DeviceID    int
	Disabled    map[Sensor]bool

	// sensors
	Temperature float64
	Humidity    uint8
	Light       uint8
	Motion      bool
	LowBattery  bool
	VCC         uint

	// private
	sensors []Sensor
}

// Parse a line produce by RF12Demo sketch
//
// Example of line generated by RF12Demo sketch that received a packet from a device:
//
//       OK 2 3 156 149 213 0
//          ^ ^ -------------
//     header |      ^
//            |  data bytes
//            |
//         node kind
//
// header:
//
//      0   0   0   0   0   0   1   0
//      ^   ^   ^   -----------------
//     CTL DST ACK         ^
//                    node id => 2
//
// node kind:
//
//      0   0   0   0   0   0   1   1
//      ^   -------------------------
// reserved            ^
//               node kind => 3
func Parse(line string) (*Packet, error) {
	log.Debugf("line: %s", line)

	result := &Packet{
		At:          time.Now(),
		Corrections: map[Sensor]float64{},
		Disabled:    map[Sensor]bool{},
	}

	// split line
	ary := strings.Split(strings.TrimSpace(line), " ")
	if (len(ary) <= 3) || (ary[0] != "OK") {
		return nil, errLineFormat
	}

	// device id
	b, err := byteFor(ary[1])
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse device id byte '%s'", ary[1])
	}

	result.DeviceID = int(b & 0x1F)

	// device kind
	b, err = byteFor(ary[2])
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse device kind byte '%s'", ary[2])
	}

	if (b & 0x80) != 0 {
		return nil, errors.New("received data with reserved field set")
	}

	result.Device = Device(int(b & 0x7F))

	// data
	data := make([]byte, len(ary)-3)
	for i, val := range ary[3:] {
		b, err = byteFor(val)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse data byte '%s'", val)
		}
		data[i] = b
	}

	expectedLen := result.Device.dataLength()
	if len(data) != expectedLen {
		return nil, fmt.Errorf("unexpected data length for packet %+v (expected %d bytes but got %d)", result, expectedLen, len(data))
	}

	if err := result.setData(data); err != nil {
		return nil, err
	}

	return result, nil
}

// ApplySettings corrects sensors values and disables sensors
func (p *Packet) ApplySettings(settings *DeviceSettings) error {
	for _, s := range settings.Sensors {
		sensor, err := SensorString(s.Name)
		if err != nil {
			return errors.Wrap(err, "unexpected sensor name")
		}

		if s.Correction != 0 {
			p.Corrections[sensor] = s.Correction

			prevValue := p.HumanizeValue(sensor)
			if err := p.correctValue(sensor, s.Correction); err != nil {
				return err
			}
			log.Debugf("Corrected %s value from '%s' to '%s'", sensor, prevValue, p.HumanizeValue(sensor))
		}

		if s.Disable {
			p.Disabled[sensor] = true
		}
	}
	return nil
}

// HasSensor returns true if packet have given sensor data
func (p *Packet) HaveSensor(sensor Sensor) bool {
	for _, s := range p.Sensors() {
		if s == sensor {
			return true
		}
	}
	return false
}

// Value returns humanized representation of sensor value with unit
func (p *Packet) HumanizeValue(sensor Sensor) string {
	switch sensor {
	case Temperature:
		return fmt.Sprintf("%.1f°", p.Temperature)
	case Humidity:
		return fmt.Sprintf("%d%%", p.Humidity)
	case Light:
		return fmt.Sprintf("%d%%", p.Light)
	case Motion:
		return fmt.Sprintf("%t", p.Motion)
	case LowBattery:
		return fmt.Sprintf("%t", p.LowBattery)
	case VCC:
		return fmt.Sprintf("%dmV", p.VCC)
	}
	return "??"
}

// Values returns humanized representations of all sensors values
func (p *Packet) HumanizeValues() []string {
	result := make([]string, len(Sensors))

	for i, sensor := range Sensors {
		if p.Device.haveSensor(sensor) {
			result[i] = p.HumanizeValue(sensor)
		}
	}

	return result
}

// HumidityStatus returns humidity status
func (p *Packet) HumidityStatus() HumidityStatus {
	if p.HaveSensor(Humidity) {
		if p.HaveSensor(Temperature) {
			if p.Temperature >= 20 && p.Temperature <= 25 &&
				p.Humidity >= 40 && p.Humidity <= 70 {
				// Comfortable if temperature is between 20-25°C and humidity between 40-70%
				return HumidityComfortable
			}
		}

		if p.Humidity < 40 {
			return HumidityDry
		}

		if p.Humidity > 70 {
			return HumidityWet
		}
	}

	return HumidityNormal
}

// Sensors returns all packet sensors
func (p *Packet) Sensors() []Sensor {
	if p.sensors == nil {
		p.sensors = []Sensor{}
		for _, sensor := range p.Device.sensors() {
			if !p.Disabled[sensor] {
				p.sensors = append(p.sensors, sensor)
			}
		}
	}
	return p.sensors
}

// String implements stringer
func (p *Packet) String() string {
	result := fmt.Sprintf("[%d][%s] ", p.DeviceID, p.Device)
	for i, sensor := range p.Device.sensors() {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprintf("%s: %s", sensor, p.HumanizeValue(sensor))
	}
	return result
}

func (p *Packet) setData(data []byte) error {
	curByte := 0
	curBytePos := 0

	for _, sensor := range p.Device.sensors() {
		var value uint64

		bits := sensor.bits()
		if bits > 64 {
			return errors.New("sensor needs more than 64 bits")
		}

		totalBitsShift := curBytePos + bits

		bytesNeeded := totalBitsShift / 8
		if (totalBitsShift % 8) != 0 {
			bytesNeeded += 1
		}

		for i := 0; i < bytesNeeded; i++ {
			value += uint64(data[curByte+i]) << uint(8*i)
		}

		value = (value >> uint(curBytePos)) & ((1 << uint(bits)) - 1)

		if err := p.setValue(sensor, value); err != nil {
			return err
		}

		curByte += (bytesNeeded - 1)
		curBytePos = totalBitsShift % 8
	}

	return nil
}

func (p *Packet) setValue(sensor Sensor, val uint64) error {
	switch sensor {
	case Temperature:
		p.Temperature = sensor.temperature(val)
	case Humidity:
		p.Humidity = sensor.humidity(val)
	case Light:
		p.Light = sensor.light(val)
	case Motion:
		p.Motion = sensor.motion(val)
	case LowBattery:
		p.LowBattery = sensor.lowBattery(val)
	case VCC:
		p.VCC = sensor.vcc(val)
	default:
		return fmt.Errorf("unexpected sensor '%s'", sensor)
	}
	return nil
}

func (p *Packet) correctValue(sensor Sensor, correction float64) error {
	switch sensor {
	case Temperature:
		p.Temperature += correction
	case Humidity:
		val := float64(p.Humidity) + correction
		if val < 0 {
			val = 0
		}
		p.Humidity = uint8(val)
	case Light:
		val := float64(p.Light) + correction
		if val < 0 {
			val = 0
		}
		p.Light = uint8(val)
	case VCC:
		val := float64(p.VCC) + correction
		if val < 0 {
			val = 0
		}
		p.VCC = uint(val)
	default:
		return fmt.Errorf("unexpected correction for sensor '%s'", sensor)
	}
	return nil
}

func byteFor(val string) (byte, error) {
	result, err := strconv.ParseUint(val, 10, 8)
	return byte(result), err
}
