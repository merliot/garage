// go run ./cmd
// go run -tags prime ./cmd
// tinygo flash -target xxx ./cmd

package main

import (
	"github.com/merliot/dean"
	"github.com/merliot/device/runner"
	"github.com/merliot/garage"
)

var (
	id           = dean.GetEnv("ID", "garage01")
	name         = dean.GetEnv("NAME", "Garage Doors")
	deployParams = dean.GetEnv("DEPLOY_PARAMS", "target=demo&door1=Car&relay1=DEMO1&trig1=DEMO2=&echo1=DEMO3")
	port         = dean.GetEnv("PORT", "8000")
	portPrime    = dean.GetEnv("PORT_PRIME", "8001")
	user         = dean.GetEnv("USER", "")
	passwd       = dean.GetEnv("PASSWD", "")
	dialURLs     = dean.GetEnv("DIAL_URLS", "")
	ssids        = dean.GetEnv("WIFI_SSIDS", "")
	passphrases  = dean.GetEnv("WIFI_PASSPHRASES", "")
)

func main() {
	garage := garage.New(id, "garage", name).(*garage.Garage)
	garage.SetDeployParams(deployParams)
	garage.SetWifiAuth(ssids, passphrases)
	garage.SetDialURLs(dialURLs)
	runner.Run(garage, port, portPrime, user, passwd, dialURLs)
}
