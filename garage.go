package garage

import (
	"net/url"
	"strconv"

	"github.com/merliot/dean"
	"github.com/merliot/device"
)

type Sonic struct {
	TrigGpio string
	EchoGpio string
	Dist     int
	Min      int
	Max      int
	targetSonicStruct
}

type Door struct {
	Name      string
	Clicked   bool
	RelayGpio string
	Sonic
	targetDoorStruct
}

type Garage struct {
	*device.Device
	Doors [2]Door
	targetStruct
}

type MsgClick struct {
	Path    string
	Door    int
	Clicked bool
}

type MsgPosition struct {
	Path string
	Door int
	Dist int
	Min  int
	Max  int
}

var targets = []string{"demo", "rpi", "nano-rp2040", "wioterminal"}

func New(id, model, name string) dean.Thinger {
	println("NEW GARAGE")
	g := &Garage{}
	g.Device = device.New(id, model, name, targets).(*device.Device)
	g.targetNew()
	return g
}

func (g *Garage) save(msg *dean.Msg) {
	msg.Unmarshal(g).Broadcast()
}

func (g *Garage) getState(msg *dean.Msg) {
	g.Path = "state"
	msg.Marshal(g).Reply()
}

func (g *Garage) click(msg *dean.Msg) {
	var msgClick MsgClick
	msg.Unmarshal(&msgClick)
	door := &g.Doors[msgClick.Door]
	door.Clicked = msgClick.Clicked
	if g.IsMetal() {
		if msgClick.Clicked {
			door.relayOn()
		}
	}
	msg.Broadcast()
}

func (g *Garage) position(msg *dean.Msg) {
	var pos MsgPosition
	msg.Unmarshal(&pos)
	msg.Unmarshal(&g.Doors[pos.Door])
	msg.Broadcast()
}

func (g *Garage) Subscribers() dean.Subscribers {
	return dean.Subscribers{
		"state":     g.save,
		"get/state": g.getState,
		"click":     g.click,
		"position":  g.position,
	}
}

func (g *Garage) setDoor(num int, name, relay, trig, echo string) {
	door := &g.Doors[num]
	door.Name = name
	door.RelayGpio = relay
	door.Sonic.TrigGpio = trig
	door.Sonic.EchoGpio = echo
}

func firstValue(values url.Values, key string) string {
	if v, ok := values[key]; ok {
		return v[0]
	}
	return ""
}

func (g *Garage) parseParams() {
	values := g.ParseDeployParams()
	for i := range g.Doors {
		num := strconv.Itoa(i + 1)
		name := firstValue(values, "door"+num)
		relay := firstValue(values, "relay"+num)
		trig := firstValue(values, "trig"+num)
		echo := firstValue(values, "echo"+num)
		g.setDoor(i, name, relay, trig, echo)
	}
}

func (g *Garage) Run(i *dean.Injector) {
	g.parseParams()
	g.run(i)
}
