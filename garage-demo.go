//go:build !tinygo && !rpi

package garage

import (
	"time"

	"github.com/merliot/dean"
)

type targetSonicStruct struct {
}

type targetDoorStruct struct {
	moving    bool
	goingDown bool
	clicked   bool
}

type targetStruct struct {
	osStruct
}

func (g *Garage) targetNew() {
	g.osNew()
}

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
			d.Dist -= 5
			if d.Dist < 0 {
				d.Dist = 0
			}
			if d.Dist == 0 {
				d.goingDown = false
				d.clicked = false
			}
		} else {
			d.Dist += 5
			if d.Dist > 100 {
				d.Dist = 100
			}
			if d.Dist == 100 {
				d.goingDown = true
				d.clicked = false
			}
		}
		d.sendPosition(inj, d.Dist)
	}
}

func (g *Garage) runDoors(inj *dean.Injector) {
	for i := range g.Doors {
		door := &g.Doors[i]
		door.run(inj)
	}
}

func (g *Garage) run(inj *dean.Injector) {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			g.runDoors(inj)
		}
	}
}
