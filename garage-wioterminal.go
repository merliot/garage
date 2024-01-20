//go:build wioterminal

package garage

import "github.com/merliot/device/target"

func (g *Garage) pins() target.GpioPins {
	return g.Targets["wioterminal"].GpioPins
}
