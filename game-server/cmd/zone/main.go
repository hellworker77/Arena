package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"game-server/internal/zone"
)

func main() {
	var listen string
	var zoneID uint
	flag.StringVar(&listen, "listen", "127.0.0.1:4000", "TCP listen address for gateway link")
	flag.UintVar(&zoneID, "zone", 1, "Zone ID")
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	s := zone.New(zone.Config{
		ListenAddr: listen,
		ZoneID:     uint32(zoneID),
		TickHz:     20,
	})
	if err := s.Start(ctx); err != nil {
		log.Fatalf("zone: %v", err)
	}
}
