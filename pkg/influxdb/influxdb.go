// Package influxdb ouputs packets to an influxdb server
package influxdb

import (
	"github.com/aymerick/homlet"
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/pkg/errors"
)

// Handler represents an influxdb handler
type Handler struct {
	URL      string // eg: http://localhost:8086
	Username string
	Password string
}

// Push sends packet to influxdb
func (h *Handler) Push(packet *homlet.Packet, settings *homlet.DeviceSettings) error {
	// Make client
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:               h.URL,
		Username:           h.Username,
		Password:           h.Password,
		InsecureSkipVerify: true,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create influxdb client")
	}
	defer c.Close()

	// // Create a new point batch
	// bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
	// 	Database:  "BumbleBeeTuna",
	// 	Precision: "s",
	// })
	//
	// // Create a point and add to batch
	// tags := map[string]string{"cpu": "cpu-total"}
	// fields := map[string]interface{}{
	// 	"idle":   10.1,
	// 	"system": 53.3,
	// 	"user":   46.6,
	// }
	// pt, err := client.NewPoint("cpu_usage", tags, fields, time.Now())
	// if err != nil {
	// 	fmt.Println("Error: ", err.Error())
	// }
	// bp.AddPoint(pt)
	//
	// // Write the batch
	// c.Write(bp)

	// FIXME
	return errors.New("not impl")
}
