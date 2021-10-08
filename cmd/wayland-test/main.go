package main

import (
	"flag"
	"log"

	"github.com/jchv/jtk/wayland"
)

func main() {
	display := flag.String("display", "", "Wayland socket to connect to, e.g. wayland-0")
	flag.Parse()

	conn, err := wayland.Connect(*display)
	if err != nil {
		log.Fatalf("Error connecting to Wayland compositor: %v", err)
	}

	defer conn.Close()
}
