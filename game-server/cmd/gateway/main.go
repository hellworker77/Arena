package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"game-server/internal/gateway"
)

func main() {
	var udpAddr string
	flag.StringVar(&udpAddr, "udp", ":7777", "UDP listen address")

	zones := make(gateway.ZoneFlags)
	flag.Var(zones, "zone", "Zone mapping: <zoneID>=<host:port> (repeatable)")

	flag.Parse()
	if len(zones) == 0 {
		log.Fatalf("provide at least one -zone (e.g. -zone 1=127.0.0.1:4000)")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	srv, err := gateway.New(gateway.Config{
		UDPListenAddr: udpAddr,
		Zones: zones,
		IdleTimeout: 30 * time.Second,
	})
	if err != nil { log.Fatalf("gateway init: %v", err) }
	if err := srv.Start(ctx); err != nil { log.Fatalf("gateway: %v", err) }
}
