package homlet

import "github.com/spf13/viper"

// DeviceSettings represents settings for a specific device
type DeviceSettings struct {
	Domoticz int    // domoticz idx
	ID       int    // device id
	Room     string // room where device is located
	Sensors  []*SensorSettings
}

// SensorSettings represents settings for a specific sensor
type SensorSettings struct {
	Name       string
	Correction float64 // correction to apply to sensor value
	Disable    bool    // disable that sensor
}

// DevicesSettings unmarshall all devices settings from conf
func DevicesSettings() ([]*DeviceSettings, error) {
	result := []*DeviceSettings{}
	if err := viper.UnmarshalKey("devices", &result); err != nil {
		return nil, err
	}
	return result, nil
}
