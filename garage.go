package garage

import (
	"fmt"
	"math"
	"net/http"

	"github.com/merliot/dean"
	"github.com/merliot/device"
	"github.com/merliot/device/relay"
	"github.com/merliot/device/vl53l1x"
)

type Sensor struct {
	vl53l1x.Vl53l1x `json:"-"`
	Dist            int32
	Min             int32
	Max             int32
	lastDist        int32
}

type Door struct {
	Name   string
	Relay  relay.Relay
	Sensor Sensor
	// for demo
	moving    bool
	goingDown bool
	clicked   bool
}

type Garage struct {
	*device.Device
	Door Door
}

type MsgClick struct {
	Path    string
	Clicked bool
}

type MsgPosition struct {
	Path string
	Dist int32
	Min  int32
	Max  int32
}

var targets = []string{"demo", "nano-rp2040"}

func New(id, model, name string) dean.Thinger {
	fmt.Println("NEW GARAGE\r")
	return &Garage{
		Device: device.New(id, model, name, fs, targets).(*device.Device),
		Door:   Door{Sensor: Sensor{Min: math.MaxInt32}},
	}
}

func (g *Garage) save(msg *dean.Msg) {
	msg.Unmarshal(g).Broadcast()
}

func (g *Garage) getState(msg *dean.Msg) {
	g.Path = "state"
	g.parseParams()
	msg.Marshal(g).Reply()
}

func (g *Garage) click(msg *dean.Msg) {
	var msgClick MsgClick
	msg.Unmarshal(&msgClick)
	g.Door.Relay.State = msgClick.Clicked
	if g.IsMetal() {
		if msgClick.Clicked {
			g.Door.relayOn()
		}
	}
	msg.Broadcast()
}

func (g *Garage) position(msg *dean.Msg) {
	msg.Unmarshal(&g.Door.Sensor).Broadcast()
}

func (g *Garage) Subscribers() dean.Subscribers {
	return dean.Subscribers{
		"state":     g.save,
		"get/state": g.getState,
		"click":     g.click,
		"position":  g.position,
	}
}

func (g *Garage) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	g.API(w, req, g)
}

func (g *Garage) parseParams() {
	g.Door.Name = g.ParamFirstValue("door")
	g.Door.Relay.Gpio = g.ParamFirstValue("relay")
}

func (g *Garage) configure() {
	g.Door.Relay.Configure()
	g.Door.Sensor.Vl53l1x.Configure()
}

func (g *Garage) Setup() {
	g.Device.Setup()
	g.parseParams()
	g.configure()
}

func (g *Garage) sendPosition(inj *dean.Injector, dist int32) {

	g.Lock()

	sensor := &g.Door.Sensor

	if dist == sensor.lastDist {
		g.Unlock()
		return
	}

	sensor.Dist = dist
	sensor.lastDist = dist

	if dist > sensor.Max {
		sensor.Max = dist
	}
	if dist < sensor.Min {
		sensor.Min = dist
	}

	var msg dean.Msg
	var pos = MsgPosition{
		Path: "position",
		Dist: sensor.Dist,
		Max:  sensor.Max,
		Min:  sensor.Min,
	}

	g.Unlock()

	inj.Inject(msg.Marshal(pos))
}
