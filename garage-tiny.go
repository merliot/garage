//go:build tinygo

package garage

import (
	"embed"
	"time"

	"github.com/merliot/dean"
)

var fs embed.FS

func (d *Door) relayOn() {
	d.Relay.On()
	time.Sleep(500 * time.Millisecond)
	d.Relay.Off()
}

func (g *Garage) Run(inj *dean.Injector) {

	// Service sensor
	sensor := &g.Door.Sensor.Vl53l1x
	for {
		if dist, ok := sensor.Distance(); ok {
			g.sendPosition(inj, dist)
		}
		time.Sleep(time.Second)
	}
}
