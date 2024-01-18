package garage

import (
	"github.com/merliot/dean"
	"github.com/merliot/relays"
)

type Garage struct {
	*relays.Relays
}

var targets = []string{"demo", "rpi", "nano-rp2040", "wioterminal"}

func New(id, model, name string) dean.Thinger {
	println("NEW GARAGE")
	g := &Garage{}
	g.Relays = relays.NewRelays(id, model, name, targets).(*relays.Relays)
	return g
}

func (g *Garage) save(msg *dean.Msg) {
	msg.Unmarshal(g).Broadcast()
}

func (g *Garage) getState(msg *dean.Msg) {
	g.Path = "state"
	msg.Marshal(g).Reply()
}

func (g *Garage) Subscribers() dean.Subscribers {
	return dean.Subscribers{
		"state":     g.save,
		"get/state": g.getState,
	}
}
