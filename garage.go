package garage

import (
	"embed"
	"math"
	"net/http"
	"net/url"

	"github.com/merliot/dean"
	"github.com/merliot/device"
)

//go:embed css html js images template
var fs embed.FS

type Sensor struct {
	Dist     int32
	Min      int32
	Max      int32
	lastDist int32
	targetSensorStruct
}

type Door struct {
	Name      string
	Clicked   bool
	RelayGpio string
	Sensor
	targetDoorStruct
}

type Garage struct {
	*device.Device
	Door Door
	targetStruct
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

var targets = []string{"demo", "nano-rp2040", "wioterminal"}

func New(id, model, name string) dean.Thinger {
	println("NEW GARAGE")
	g := &Garage{}
	g.Device = device.New(id, model, name, fs, targets).(*device.Device)
	g.Door.Min = math.MaxInt32
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
	g.Door.Clicked = msgClick.Clicked
	if g.IsMetal() {
		if msgClick.Clicked {
			g.Door.relayOn()
		}
	}
	msg.Broadcast()
}

func (g *Garage) position(msg *dean.Msg) {
	msg.Unmarshal(&g.Door).Broadcast()
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

func firstValue(values url.Values, key string) string {
	if v, ok := values[key]; ok {
		return v[0]
	}
	return ""
}

func (g *Garage) parseParams() {
	values := g.ParseDeployParams()
	g.Door.Name = firstValue(values, "door")
	g.Door.RelayGpio = firstValue(values, "relay")
}

func (g *Garage) Run(i *dean.Injector) {
	g.parseParams()
	g.run(i)
}

func (d *Door) sendPosition(inj *dean.Injector, dist int32) {

	if dist == d.lastDist {
		return
	}

	d.Dist = dist
	d.lastDist = dist

	if dist > d.Max {
		d.Max = dist
	}
	if dist < d.Min {
		d.Min = dist
	}

	var msg dean.Msg
	var pos = MsgPosition{
		Path: "position",
		Dist: d.Dist,
		Max:  d.Max,
		Min:  d.Min,
	}

	inj.Inject(msg.Marshal(pos))
}
