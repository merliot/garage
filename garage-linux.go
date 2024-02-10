//go:build !tinygo

package garage

import (
	"net/http"
)

func (g *Garage) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	g.API(w, req, g)
}

func (g *Garage) Icon() []byte {
	icon, _ := fs.ReadFile("images/icon.png")
	return icon
}

func (g *Garage) DescHtml() []byte {
	desc, _ := fs.ReadFile("html/desc.html")
	return desc
}

func (g *Garage) SupportedTargets() string {
	return g.Targets.FullNames()
}
