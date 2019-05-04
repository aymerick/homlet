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

		log.Infof("Packet received: %s", packet)

		// send to domoticz
		if s.domoticz != nil {
			s.domoticz.Push(packet, s.deviceSettings(packet.DeviceID))
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
