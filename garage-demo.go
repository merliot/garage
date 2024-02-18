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
func (d *Door) run(inj *dean.Injector) {
	sensor := &d.Sensor
	if d.clicked {
		if d.goingDown {
			sensor.Dist -= 5
			if sensor.Dist < 0 {
				sensor.Dist = 0
			}
			if sensor.Dist == 0 {
				d.goingDown = false
				d.clicked = false
			}
		} else {
			sensor.Dist += 5
			if sensor.Dist > 100 {
				sensor.Dist = 100
			}
			if sensor.Dist == 100 {
				d.goingDown = true
				d.clicked = false
			}
		}
		d.sendPosition(inj, sensor.Dist)
	}
}

func (g *Garage) Run(inj *dean.Injector) {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			g.Door.run(inj)
		}
	}
}
