package homlet

import (
	"sync"
)

//go:generate enumer -type=Device -output=device_gen.go -transform=snake

// Device represents a device kind
type Device int

// Supported devices
const (
	Unknown     Device = iota
	JeenodeTHLM        // Jeenode: Temperature Humidity Light Motion
	JeenodeTHL         // Jeenode: Temperature Humidity Light
	TinytxT            //  TinyTX: Temperature
	TinytxTH           //  TinyTX: Temperature Humidity
	TinytxTL           //  TinyTX: Temperature Light
)

var deviceSensors = map[Device][]Sensor{
	JeenodeTHLM: {Temperature, Humidity, Light, Motion, LowBattery},
	JeenodeTHL:  {Temperature, Humidity, Light, LowBattery},
	TinytxT:     {Temperature, VCC},
	TinytxTH:    {Temperature, Humidity, VCC},
	TinytxTL:    {Temperature, Light, VCC},
}

var deviceDataLength map[Device]int
var once sync.Once

func computeDeviceDataLength() map[Device]int {
	result := map[Device]int{}
	for device, sensors := range deviceSensors {
		bits := 0
		for _, sensor := range sensors {
			bits += sensor.bits()
		}
		result[device] = bits / 8
		if (bits % 8) != 0 {
			result[device] += 1
		}
	}
	return result
}

func (d Device) dataLength() int {
	once.Do(func() {
		deviceDataLength = computeDeviceDataLength()
	})
	return deviceDataLength[d]
}

func (d Device) haveSensor(sensor Sensor) bool {
	for _, s := range d.sensors() {
		if s == sensor {
			return true
		}
	}
	return false
}

func (d Device) sensors() []Sensor {
	return deviceSensors[d]
}
