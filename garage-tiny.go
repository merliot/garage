//go:build tinygo

package garage

import (
	"machine"
	"time"

	"github.com/merliot/dean"
	"tinygo.org/x/drivers/hcsr04"
)

type targetSonicStruct struct {
	hcsr04    hcsr04.Device
	hasHcsr04 bool
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

	// Configure doors
	for i := range g.Doors {
		door := &g.Doors[i]
		// Configure relay
		door.relayPin = machine.NoPin
		if pin, ok := g.Pins()[door.RelayGpio]; ok {
			door.relayPin = machine.Pin(pin)
			door.relayPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
		}
		// Configure HC-SR04 ultrasonic distance sensor
		if pin, ok := g.Pins()[door.TrigGpio]; ok {
			trigPin := machine.Pin(pin)
			if pin, ok := g.Pins()[door.EchoGpio]; ok {
				echoPin := machine.Pin(pin)
				door.hasHcsr04 = true
				door.hcsr04 = hcsr04.New(trigPin, echoPin)
				door.hcsr04.Configure()
			}
		}
	}

	// Service HC-SR04 sensors
	for {
		for i := range g.Doors {
			door := &g.Doors[i]
			if door.hasHcsr04 {
				dist := door.hcsr04.ReadDistance() // mm
				door.sendPosition(inj, dist)
			}
			time.Sleep(500 * time.Millisecond)
		}
	}
}
