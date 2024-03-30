//go:build !tinygo && !rpi

package garage

import (
	"time"

	"github.com/merliot/dean"
)

func (d *Door) relayOn() {
	if d.clicked {
		// reverse direction if stopping
		d.goingDown = !d.goingDown
	}
	d.clicked = !d.clicked
}

// Simulate a garage door
func (g *Garage) run(inj *dean.Injector) {
	door := &g.Door
	sensor := &g.Door.Sensor
	if door.clicked {
		if door.goingDown {
			sensor.Dist -= 5
			if sensor.Dist < 0 {
				sensor.Dist = 0
			}
			if sensor.Dist == 0 {
				door.goingDown = false
				door.clicked = false
			}
		} else {
			sensor.Dist += 5
			if sensor.Dist > 100 {
				sensor.Dist = 100
			}
			if sensor.Dist == 100 {
				door.goingDown = true
				door.clicked = false
			}
		}
		g.sendPosition(inj, sensor.Dist)
	}
}

func (g *Garage) Run(inj *dean.Injector) {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			g.run(inj)
		}
	}
}
