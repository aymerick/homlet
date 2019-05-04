// Package domoticz implements a domoticz output
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

	if !packet.Device.HaveSensor(homlet.Temperature) && !packet.Device.HaveSensor(homlet.Humidity) {
		return nil
	}

	url := fmt.Sprintf("%s/json.htm?type=command&param=udevice&idx=%d&%s", h.URL, settings.Domoticz, h.paramsForPacket(packet))

	log.Debugf("Pushing to domoticz: %s", url)

	// FIXME check if domoticz support a POST instead
	resp, err := http.Get(url)
	if err != nil {
		log.WithError(err).Error("Failed to push packet to domoticz")
		return err
	}

	// check response
	respText, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Warn("Failed to get domoticz response")
	} else {
		log.Info("Packet pushed to domoticz")
		log.Debugf("Domoticz response: %s", respText)
	}

	return nil
}

func (h *Handler) paramsForPacket(packet *homlet.Packet) string {
	if packet.Device.HaveSensor(homlet.Temperature) && packet.Device.HaveSensor(homlet.Humidity) {
		return fmt.Sprintf("nvalue=0&svalue=%.1f;%d;0", packet.Temperature, packet.Humidity)
	} else if packet.Device.HaveSensor(homlet.Temperature) {
		return fmt.Sprintf("nvalue=0&svalue=%.1f", packet.Temperature)
	} else if packet.Device.HaveSensor(homlet.Humidity) {
		return fmt.Sprintf("nvalue=%d&svalue=0", packet.Humidity)
	}
	return ""
}
