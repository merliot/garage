package main

import (
	"log"
	"os"

	"github.com/merliot/garage"
)

//go:generate go run main.go
func main() {
	garage := garage.New("proto", "garage", "proto").(*garage.Garage)
	if err := garage.GenerateUf2s("../.."); err != nil {
		log.Println("Error generating UF2s:", err)
		os.Exit(1)
	}
}
