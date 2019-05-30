package server

import (
	"github.com/aymerick/homlet"
	"github.com/aymerick/homlet/pkg/domoticz"
	log "github.com/sirupsen/logrus"
)

// Server receives packets
type Server struct {
	domoticz *domoticz.Handler
	packets  chan *homlet.Packet
	settings []*homlet.DeviceSettings
}

// New instanciates a new server
func New(c chan *homlet.Packet, settings []*homlet.DeviceSettings) *Server {
	return &Server{
		packets:  c,
		settings: settings,
	}
}

// SetDomoticz sets domoticz handler
func (s *Server) SetDomoticz(domo *domoticz.Handler) {
	s.domoticz = domo
}

// Run blocks until server is stopped
func (s *Server) Run() {
	for packet := range s.packets {
		if packet == nil {
			// packet reader ended
			return
		}

		settings := s.deviceSettings(packet.DeviceID)

		// log
		log.Infof("Received: [%s]%s", settings.Room, packet)

		// correct values and disable sensors
		if err := packet.ApplySettings(settings); err != nil {
			log.WithError(err).Error("Failed to set settings to packet")
			continue
		}

		// send to domoticz
		if (s.domoticz != nil) && (settings != nil) {
			if err := s.domoticz.Push(packet, settings); err != nil {
				log.WithError(err).Error("Failed to push packet to domoticz")
			}
		}

		// FIXME send to influxdb
	}
}

func (s *Server) deviceSettings(id int) *homlet.DeviceSettings {
	for _, s := range s.settings {
		if s.ID == id {
			return s
		}
	}
	return nil
}
