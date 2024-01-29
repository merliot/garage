//go:build tinygo

package garage

import (
	"machine"
	"time"

	"github.com/merliot/dean"
	"tinygo.org/x/drivers/vl53l1x"
)

type targetSensorStruct struct {
	vl53l1x    vl53l1x.Device
	hasVl53l1x bool
}

type targetDoorStruct struct {
	relayPin machine.Pin
}

type targetStruct struct {
}

func (g *Garage) targetNew() {
}

func (d *Door) relayOn() {
	if d.relayPin != machine.NoPin {
		d.relayPin.High()
		time.Sleep(500 * time.Millisecond)
		d.relayPin.Low()
	}
}

func (g *Garage) run(inj *dean.Injector) {

	// Configure relay
	door := &g.Door
	door.relayPin = machine.NoPin
	if pin, ok := g.Pins()[door.RelayGpio]; ok {
		door.relayPin = machine.Pin(pin)
		door.relayPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	}

	// Configure VL53L1x time-of-flight distance sensor
	machine.I2C0.Configure(machine.I2CConfig{
		Frequency: 400000,
	})
	door.vl53l1x = vl53l1x.New(machine.I2C0)
	door.hasVl53l1x = door.vl53l1x.Connected()
	if door.hasVl53l1x {
		door.vl53l1x.Configure(true)
		door.vl53l1x.SetMeasurementTimingBudget(50000)
		door.vl53l1x.StartContinuous(50)
	}

	// Service sensor
	for {
		door := &g.Door
		if door.hasVl53l1x {
			door.vl53l1x.Read(true)
			dist := door.vl53l1x.Distance() // mm
			door.sendPosition(inj, dist)
		}
		time.Sleep(time.Second)
	}
}
