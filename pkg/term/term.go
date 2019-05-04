// Package term implements terminal UI
package term

import (
	"fmt"
	"strings"
	"time"

	"github.com/aymerick/homlet"
	"github.com/dustin/go-humanize"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// UI is the terminal UI
type UI struct {
	devices  map[int]*homlet.Packet
	header   []string
	packets  chan *homlet.Packet
	settings []*homlet.DeviceSettings
	table    *widgets.Table
}

// NewUI instanciates a new terminal UI
func NewUI(c chan *homlet.Packet, settings []*homlet.DeviceSettings) *UI {
	ui := &UI{
		devices:  map[int]*homlet.Packet{},
		header:   makeHeader(),
		packets:  c,
		settings: settings,
		table:    makeTable(),
	}

	ui.render()

	return ui
}

// Run displays UI, and blocks until user wants to quit
func (ui *UI) Run() {
	events := termui.PollEvents()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case packet := <-ui.packets:
			if packet == nil {
				// packet reader ended
				return
			}

			ui.devices[packet.DeviceID] = packet
			ui.render()

		case e := <-events:
			switch e.ID {
			case "q", "<C-c>":
				// quit
				return
			}

		case <-ticker.C:
			ui.render()
		}
	}
}

func (ui *UI) render() {
	ui.table.Rows = ui.rows()
	termui.Render(ui.table)
}

func (ui *UI) rows() [][]string {
	result := make([][]string, len(ui.devices)+1)
	result[0] = ui.header

	// sort packets
	packets := make([]*homlet.Packet, len(ui.devices))
	i := 0
	for _, packet := range ui.devices {
		packets[i] = packet
		i++
	}
	homlet.SortPackets(packets)

	// build rows
	for i, packet := range packets {
		result[i+1] = ui.makeFields(packet)
	}

	return result
}

func (ui *UI) deviceSettings(id int) *homlet.DeviceSettings {
	for _, s := range ui.settings {
		if s.ID == id {
			return s
		}
	}
	return nil
}

func (ui *UI) makeFields(packet *homlet.Packet) []string {
	settings := ui.deviceSettings(packet.DeviceID)
	room := ""
	if settings != nil {
		room = settings.Room
	}

	values := packet.Values()
	result := make([]string, len(values)+4)
	result[0] = fmt.Sprintf("%d", packet.DeviceID)
	result[1] = room
	result[2] = packet.Device.String()
	copy(result[3:], values)
	result[len(result)-1] = humanize.Time(packet.At)
	return result
}

func makeHeader() []string {
	result := make([]string, len(homlet.Sensors)+4)
	result[0] = "ID"
	result[1] = "ROOM"
	result[2] = "KIND"
	for i, sensor := range homlet.Sensors {
		result[i+3] = strings.ToUpper(sensor.String())
	}
	result[len(result)-1] = "LAST SEEN"
	return result
}

func makeTable() *widgets.Table {
	result := widgets.NewTable()
	result.TextStyle = termui.NewStyle(termui.ColorWhite)
	result.RowSeparator = true
	result.TextAlignment = termui.AlignCenter
	result.FillRow = true
	result.ColumnWidths = []int{10, 20, 20, 20, 10, 10, 10, 20, 10, 30}

	result.BorderStyle = termui.NewStyle(termui.ColorGreen)
	result.SetRect(0, 0, 170, 50)
	return result
}
