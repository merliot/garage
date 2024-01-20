//go:build nano_rp2040

package garage

import "github.com/merliot/device/target"

func (g *Garage) pins() target.GpioPins {
	return g.Targets["nano-rp2040"].GpioPins
}
