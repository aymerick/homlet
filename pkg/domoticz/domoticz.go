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
	hdwID   int
	motions map[int]bool
	url     string
}

// New instanciates a new domoticz handler
func New(hdwID int, url string) *Handler {
	return &Handler{
		hdwID:   hdwID,
		motions: map[int]bool{},
		url:     url,
	}
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

	// command
	url := fmt.Sprintf("%s/json.htm?type=command&%s", h.url, p)

	// battery level
	if bat := h.batteryLevel(packet); bat != "" {
		url += fmt.Sprintf("&%s", bat)
	}

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

// https://www.domoticz.com/wiki/Domoticz_API/JSON_URL's#Additional_parameters_.28signal_level_.26_battery_level.29
func (h *Handler) batteryLevel(packet *homlet.Packet) string {
	if packet.HaveSensor(homlet.VCC) {
		// NOTE that:
		//   3800mv is max allowed voltage
		//   2200mv is min voltage for RF12B chip

		// 3300mv => 100%
		// 2300mv =>   0%
		val := (packet.VCC - 2300) / 10
		if packet.VCC < 2300 {
			val = 0
		} else if packet.VCC > 3300 {
			val = 100
		}
		return fmt.Sprintf("battery=%d", val)
	} else if packet.HaveSensor(homlet.LowBattery) {
		val := 100
		if packet.LowBattery {
			// if LowBattery is true then it means that voltage dropped under 3100mv
			val = 5 // FIXME arbitrary value
		}
		return fmt.Sprintf("battery=%d", val)
	}

	// not supported
	return ""
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
		prev, alreadySent := h.motions[packet.DeviceID]

		// send to domoticz only if state changed
		if !alreadySent || (packet.Motion != prev) {
			// https://www.domoticz.com/wiki/Domoticz_API/JSON_URL's#Turn_a_light.2Fswitch_on
			cmd := "Off"
			if packet.Motion {
				cmd = "On"
			}

			h.motions[packet.DeviceID] = packet.Motion
			return fmt.Sprintf("param=switchlight&idx=%d&switchcmd=%s", deviceID, cmd)
		}

		// state did not changed
		return ""
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
