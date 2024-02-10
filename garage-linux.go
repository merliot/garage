//go:build !tinygo

package garage

import (
	"net/http"
)

func (g *Garage) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	g.API(w, req, g)
}
