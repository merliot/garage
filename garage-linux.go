//go:build !tinygo

package garage

import (
	"embed"
	"html/template"
	"net/http"
	"strings"

	"github.com/merliot/device"
)

//go:embed css js images template
var fs embed.FS

type osStruct struct {
	templates *template.Template
}

func (g *Garage) osNew() {
	g.CompositeFs.AddFS(fs)
	g.templates = g.CompositeFs.ParseFS("template/*")
}

func (g *Garage) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch strings.TrimPrefix(req.URL.Path, "/") {
	case "state":
		device.ShowState(g.templates, w, g)
	default:
		g.API(g.templates, w, req)
	}
}
