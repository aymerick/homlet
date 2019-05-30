// Package domoticz outputs packets to a domoticz instance
//
// Documentation:
//   https://www.domoticz.com/wiki/Domoticz_API/JSON_URL's
package domoticz

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aymerick/homlet"
	log "github.com/sirupsen/logrus"
)

// Handler represents a domoticz handler
type Handler struct {
	HardwareId int
	URL        string
}

// Push sends packet to domoticz
func (h *Handler) Push(packet *homlet.Packet, settings *homlet.DeviceSettings) error {
	for _, sensor := range packet.Sensors() {
		sensorSettings, err := settings.Sensor(sensor)
		if err != nil {
			return err
		}

		if (sensorSettings == nil) || (sensorSettings.Domoticz == 0) {
			continue
		}

		if (sensor == homlet.Humidity) && packet.HaveSensor(homlet.Temperature) {
			// we send a single command for Temp+Humi sensors
			continue
		}

		if err := h.pushCommand(packet, sensor, sensorSettings); err != nil {
			log.WithError(err).Errorf("Failed to push packet to domoticz for sensor '%s'", sensor.String())
		}
	}
	return nil
}

func (h *Handler) pushCommand(packet *homlet.Packet, sensor homlet.Sensor, settings *homlet.SensorSettings) error {
	p := h.params(packet, sensor, settings.Domoticz)
	if p == "" {
		return nil
	}

	// FIXME ?? send vcc as battery level: https://www.domoticz.com/wiki/Domoticz_API/JSON_URL's#Additional_parameters_.28signal_level_.26_battery_level.29
	url := fmt.Sprintf("%s/json.htm?type=command&%s", h.URL, p)

	log.Debugf("Pushing to domoticz: %s", url)

	// FIXME check if domoticz support a POST instead
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	// check response
	respText, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Warn("Failed to get domoticz response")
	} else {
		log.Debugf("Domoticz response: %s", respText)
	}

	return nil
}

func (h *Handler) params(packet *homlet.Packet, sensor homlet.Sensor, deviceID int) string {
	switch sensor {
	case homlet.Temperature:
		if packet.HaveSensor(homlet.Humidity) {
			// https://www.domoticz.com/wiki/Domoticz_API/JSON_URL's#Temperature.2Fhumidity
			return fmt.Sprintf("param=udevice&idx=%d&nvalue=0&svalue=%.1f;%d;%d", deviceID, packet.Temperature, packet.Humidity, humidityStatus(packet.HumidityStatus()))
		}
		// https://www.domoticz.com/wiki/Domoticz_API/JSON_URL's#Temperature
		return fmt.Sprintf("param=udevice&idx=%d&nvalue=0&svalue=%.1f", deviceID, packet.Temperature)

	case homlet.Humidity:
		// https://www.domoticz.com/wiki/Domoticz_API/JSON_URL's#Humidity
		return fmt.Sprintf("param=udevice&idx=%d&nvalue=%d&svalue=%d", deviceID, packet.Humidity, humidityStatus(packet.HumidityStatus()))

	case homlet.Light:
		// https://www.domoticz.com/wiki/Domoticz_API/JSON_URL's#Percentage
		return fmt.Sprintf("param=udevice&idx=%d&nvalue=0&svalue=%d", deviceID, packet.Light)

	case homlet.Motion:
		// https://www.domoticz.com/wiki/Domoticz_API/JSON_URL's#Turn_a_light.2Fswitch_on
		cmd := "Off"
		if packet.Motion {
			cmd = "On"
		}
		return fmt.Sprintf("param=switchlight&idx=%d&switchcmd=%s", deviceID, cmd)
	}

	// unsupported sensor
	return ""
}

func humidityStatus(status homlet.HumidityStatus) int {
	switch status {
	case homlet.HumidityComfortable:
		return 1
	case homlet.HumidityDry:
		return 2
	case homlet.HumidityWet:
		return 3
	}
	// Normal
	return 0
}
