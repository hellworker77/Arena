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
	var proto uint
	flag.StringVar(&udpAddr, "udp", ":7777", "UDP listen address")
	flag.UintVar(&proto, "proto", 1, "UDP protocol version")
	var rateBps int
	var burst int
	var maxRel int
	flag.IntVar(&rateBps, "rateBps", 20000, "per-session UDP send rate bytes/sec")
	flag.IntVar(&burst, "burst", 40000, "per-session burst bytes")
	flag.IntVar(&maxRel, "maxReliableBytes", 65536, "max pending reliable bytes per session")

	zones := make(gateway.ZoneFlags)
	flag.Var(zones, "zone", "Zone mapping: <zoneID>=<host:port> (repeatable)")
	flag.Parse()

	if len(zones) == 0 {
		log.Fatalf("provide at least one -zone")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	srv, err := gateway.New(gateway.Config{
		UDPListenAddr: udpAddr,
		Zones: zones,
		IdleTimeout: 30*time.Second,
		ProtoVersion: uint16(proto),
		TransferTimeout: 3*time.Second,
		RateBytesPerSec: rateBps,
		BurstBytes: burst,
		MaxReliableBytes: maxRel,
	})
	if err != nil { log.Fatalf("gateway init: %v", err) }
	if err := srv.Start(ctx); err != nil { log.Fatalf("gateway: %v", err) }
}
