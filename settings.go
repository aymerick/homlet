package homlet

import "github.com/pkg/errors"

// DeviceSettings represents settings for a specific device
type DeviceSettings struct {
	ID      int    // device id
	Room    string // room where device is located
	Sensors []*SensorSettings
}

// SensorSettings represents settings for a specific sensor
type SensorSettings struct {
	Name       string
	Correction float64 // correction to apply to sensor value
	Disable    bool    // disable that sensor
	Domoticz   int     // domoticz idx
}

// Sensor returns sensor settings
func (ds *DeviceSettings) Sensor(sensor Sensor) (*SensorSettings, error) {
	for _, sensorSettings := range ds.Sensors {
		s, err := SensorString(sensorSettings.Name)
		if err != nil {
			return nil, errors.Wrap(err, "unexpected sensor name")
		}

		if s == sensor {
			return sensorSettings, nil
		}
	}

	// not found
	return nil, nil
}
