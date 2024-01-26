//go:build tinygo

package garage

import (
	"machine"
	"time"

	"github.com/merliot/dean"
	"tinygo.org/x/drivers/hcsr04"
	"tinygo.org/x/drivers/vl53l1x"
)

type targetSonicStruct struct {
	hcsr04     hcsr04.Device
	vl53l1x    vl53l1x.Device
	hasHcsr04  bool
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
		// Configure VL53L1x time-of-flight distance sensor
		if i == 0 {
			machine.I2C0.Configure(machine.I2CConfig{
				Frequency: 400000,
			})
			door.vl53l1x = vl53l1x.New(machine.I2C0)
		} else {
			machine.I2C1.Configure(machine.I2CConfig{
				Frequency: 400000,
			})
			door.vl53l1x = vl53l1x.New(machine.I2C1)
		}
		door.hasVl53l1x = door.vl53l1x.Connected()
		if door.hasVl53l1x {
			door.hasHcsr04 = false
			door.vl53l1x.Configure(true)
			door.vl53l1x.SetMeasurementTimingBudget(50000)
			door.vl53l1x.StartContinuous(50)
		}
	}

	// Service sensors
	for {
		for i := range g.Doors {
			door := &g.Doors[i]
			switch {
			case door.hasHcsr04:
				dist := door.hcsr04.ReadDistance() // mm
				door.sendPosition(inj, dist)
			case door.hasVl53l1x:
				door.vl53l1x.Read(true)
				dist := door.vl53l1x.Distance() // mm
				door.sendPosition(inj, dist)
			}
			time.Sleep(500 * time.Millisecond)
		}
	}
}
