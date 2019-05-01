// Code generated by "enumer -type=Sensor -output=sensor_gen.go -transform=snake"; DO NOT EDIT.

package homlet

import (
	"fmt"
)

const _SensorName = "temperaturehumiditylightmotionlow_batteryvcc"

var _SensorIndex = [...]uint8{0, 11, 19, 24, 30, 41, 44}

func (i Sensor) String() string {
	i -= 1
	if i < 0 || i >= Sensor(len(_SensorIndex)-1) {
		return fmt.Sprintf("Sensor(%d)", i+1)
	}
	return _SensorName[_SensorIndex[i]:_SensorIndex[i+1]]
}

var _SensorValues = []Sensor{1, 2, 3, 4, 5, 6}

var _SensorNameToValueMap = map[string]Sensor{
	_SensorName[0:11]:  1,
	_SensorName[11:19]: 2,
	_SensorName[19:24]: 3,
	_SensorName[24:30]: 4,
	_SensorName[30:41]: 5,
	_SensorName[41:44]: 6,
}

// SensorString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func SensorString(s string) (Sensor, error) {
	if val, ok := _SensorNameToValueMap[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to Sensor values", s)
}

// SensorValues returns all values of the enum
func SensorValues() []Sensor {
	return _SensorValues
}

// IsASensor returns "true" if the value is listed in the enum definition. "false" otherwise
func (i Sensor) IsASensor() bool {
	for _, v := range _SensorValues {
		if i == v {
			return true
		}
	}
	return false
}