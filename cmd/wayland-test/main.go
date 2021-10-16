package main

import (
	"flag"
	"log"

	"github.com/jchv/jtk/internal/wayland"
)

func main() {
	display := flag.String("display", "", "Wayland socket to connect to, e.g. wayland-0")
	flag.Parse()

	conn, err := wayland.Connect(*display)
	if err != nil {
		log.Fatalf("Error connecting to Wayland compositor: %v", err)
	}

	cb, err := conn.Sync()
	if err != nil {
		log.Fatalf("Error getting synchronization checkpoint: %v", err)
	}

	conn.RegisterHandler(cb.ID(), wayland.HandlerFunc(func(event wayland.Event) {
		switch t := event.(type) {
		case *wayland.WlCallbackDoneEvent:
			log.Printf("Synchronized! %d", t.CallbackData)
			conn.Close()
		}
	}))

	registry, err := conn.Registry()
	if err != nil {
		log.Fatalf("Error getting Wayland registry: %v", err)
	}

	log.Printf("registry ID: %d\n", registry.ID())

	if err := conn.EventLoop(); err != nil {
		log.Printf("Error in event loop: %v", err)
	}

	defer conn.Close()
}
