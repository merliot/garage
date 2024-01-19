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

func (g *Garage) sendPosition(i *dean.Injector, index int) {
	var msg dean.Msg
	var d = &g.Doors[index]
	var pos = MsgPosition{
		Path: "position",
		Door: index,
		Dist: d.Dist,
		Max:  d.Max,
		Min:  d.Min,
	}
	i.Inject(msg.Marshal(pos))
}

func (g *Garage) runDoor(i *dean.Injector, index int) {
	d := &g.Doors[index]
	if d.clicked {
		if d.goingDown {
			d.Dist -= 5
			if d.Dist < 0 {
				d.Dist = 0
			}
			if d.Dist < d.Min {
				d.Min = d.Dist
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
			if d.Dist > d.Max {
				d.Max = d.Dist
			}
			if d.Dist == 100 {
				d.goingDown = true
				d.clicked = false
			}
		}
		g.sendPosition(i, index)
	}
}

func (g *Garage) runDoors(i *dean.Injector) {
	for index := range g.Doors {
		g.runDoor(i, index)
	}
}

func (g *Garage) run(i *dean.Injector) {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			g.runDoors(i)
		}
	}
}
