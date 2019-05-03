package homlet

//go:generate enumer -type=Sensor -output=sensor_gen.go -transform=snake

// Sensor represents a sensor kind
type Sensor int

// Supported sensors
const (
	Temperature Sensor = iota + 1
	Humidity
	Light
	Motion
	LowBattery
	VCC // Supply voltage
)

// Sensors holds all sensors
var Sensors = []Sensor{
	Temperature, Humidity, Light, Motion, LowBattery, VCC,
}

var sensorBits = map[Sensor]int{
	Temperature: 10, // [10 bits] Temperature: -512..+512 (tenths)
	Humidity:    7,  //  [7 bits] Humidity: 0..100
	Light:       8,  //  [8 bits] Light: 0..255
	Motion:      1,  //   [1 bit] Motion: 0..1
	LowBattery:  1,  //   [1 bit] Low Battery: 0..1
	VCC:         12, // [12 bits] Supply voltage: 0..4095 mV
}

func (s Sensor) bits() int {
	return sensorBits[s]
}

func (s Sensor) temperature(value uint64) float64 {
	result := int64(value)

	if result > 512 {
		// negative value
		result -= 1024
	}

	return float64(result) / 10
}

func (s Sensor) humidity(value uint64) uint8 {
	return uint8(value)
}

func (s Sensor) light(value uint64) uint8 {
	return uint8((value * 100) / 255)
}

func (s Sensor) motion(value uint64) bool {
	return (value != 0)
}

func (s Sensor) lowBattery(value uint64) bool {
	return (value != 0)
}

func (s Sensor) vcc(value uint64) uint {
	return uint(value)
}
