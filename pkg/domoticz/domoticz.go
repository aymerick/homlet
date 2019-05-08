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
	if settings == nil || settings.Domoticz == 0 {
		return nil
	}

	if !packet.HaveSensor(homlet.Temperature) && !packet.HaveSensor(homlet.Humidity) {
		return nil
	}

	url := fmt.Sprintf("%s/json.htm?type=command&param=udevice&idx=%d&%s", h.URL, settings.Domoticz, h.paramsForPacket(packet))

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

func (h *Handler) paramsForPacket(packet *homlet.Packet) string {
	hasTemp := packet.HaveSensor(homlet.Temperature)
	hasHumi := packet.HaveSensor(homlet.Humidity)

	if hasTemp && hasHumi {
		return fmt.Sprintf("nvalue=0&svalue=%.1f;%d;%d", packet.Temperature, packet.Humidity, humidityStatus(packet.HumidityStatus()))
	} else if hasTemp {
		return fmt.Sprintf("nvalue=0&svalue=%.1f", packet.Temperature)
	} else if hasHumi {
		return fmt.Sprintf("nvalue=%d&svalue=%d", packet.Humidity, humidityStatus(packet.HumidityStatus()))
	}
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
