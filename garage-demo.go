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
	if d.clicked {
		if d.goingDown {
			d.Sensor.Dist -= 5
			if d.Sensor.Dist < 0 {
				d.Sensor.Dist = 0
			}
			if d.Sensor.Dist == 0 {
				d.goingDown = false
				d.clicked = false
			}
		} else {
			d.Sensor.Dist += 5
			if d.Sensor.Dist > 100 {
				d.Sensor.Dist = 100
			}
			if d.Sensor.Dist == 100 {
				d.goingDown = true
				d.clicked = false
			}
		}
		d.sendPosition(inj, d.Sensor.Dist)
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
