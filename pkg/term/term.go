// Package term implements terminal UI
package term

import (
	"fmt"
	"strings"
	"time"

	"github.com/aymerick/homlet"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// UI is the terminal UI
type UI struct {
	devices map[int][]string
	header  []string
	packets chan *homlet.Packet
	table   *widgets.Table
}

// NewUI instanciates a new terminal UI
func NewUI(c chan *homlet.Packet) *UI {
	result := &UI{
		devices: map[int][]string{},
		header:  makeHeader(),
		packets: c,
		table:   makeTable(),
	}

	// render table
	result.table.Rows = result.rows()
	termui.Render(result.table)

	return result
}

// Run displays UI, and blocks until user wants to quit
func (ui *UI) Run() {
	events := termui.PollEvents()

	for {
		select {
		case packet := <-ui.packets:
			if packet == nil {
				return
			}

			ui.devices[packet.DeviceID] = makeFields(packet)

			// re-render
			ui.table.Rows = ui.rows()
			termui.Render(ui.table)

		case e := <-events:
			switch e.ID {
			case "q", "<C-c>":
				// quit
				return
			}
		}
	}
}

func (ui *UI) rows() [][]string {
	result := make([][]string, len(ui.devices)+1)
	result[0] = ui.header

	// FIXME order by device id
	i := 1
	for _, device := range ui.devices {
		result[i] = device
		i++
	}

	return result
}

func makeTable() *widgets.Table {
	result := widgets.NewTable()
	result.TextStyle = termui.NewStyle(termui.ColorWhite)
	result.RowSeparator = true
	result.TextAlignment = termui.AlignCenter
	result.FillRow = true
	result.ColumnWidths = []int{10, 20, 20, 10, 10, 10, 20, 10, 30}

	result.BorderStyle = termui.NewStyle(termui.ColorGreen)
	result.SetRect(0, 0, 150, 50)
	return result
}

func makeHeader() []string {
	result := make([]string, len(homlet.Sensors)+3)
	result[0] = "ID"
	result[1] = "KIND"
	for i, sensor := range homlet.Sensors {
		result[i+2] = strings.ToUpper(sensor.String())
	}
	result[len(result)-1] = "LAST SEEN"
	return result
}

func makeFields(packet *homlet.Packet) []string {
	values := packet.Values()
	result := make([]string, len(values)+3)
	result[0] = fmt.Sprintf("%d", packet.DeviceID)
	result[1] = packet.Device.String()
	copy(result[2:], values)
	result[len(result)-1] = time.Now().Format("2006-01-02T15:04:05")
	return result
}
