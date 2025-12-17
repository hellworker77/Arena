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
	var zoneAddr string
	flag.StringVar(&udpAddr, "udp", ":7777", "UDP listen address")
	flag.StringVar(&zoneAddr, "zone", "127.0.0.1:4000", "Zone TCP address")
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	srv, err := gateway.New(gateway.Config{
		UDPListenAddr: udpAddr,
		ZoneTCPAddr:   zoneAddr,
		IdleTimeout:   30 * time.Second,
	})
	if err != nil {
		log.Fatalf("gateway init: %v", err)
	}

	if err := srv.Start(ctx); err != nil {
		log.Fatalf("gateway: %v", err)
	}
}
