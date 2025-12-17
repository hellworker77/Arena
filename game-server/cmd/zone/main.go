package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"game-server/internal/persist"
	"game-server/internal/zone"
)

func main() {
	var listen string
	var zoneID uint
	var targetZone uint
	var boundary int
	var storeDir string

	flag.StringVar(&listen, "listen", "127.0.0.1:4000", "TCP listen address for gateway link")
	flag.UintVar(&zoneID, "zone", 1, "Zone ID")
	flag.UintVar(&targetZone, "targetZone", 2, "Target zone ID for transfers")
	flag.IntVar(&boundary, "xferBoundary", 100, "Transfer boundary on X (>, or < if negative)")
	flag.StringVar(&storeDir, "store", "./data", "store directory (JSON placeholder)")
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	store, err := persist.NewJSONStore(storeDir)
	if err != nil { log.Fatalf("store: %v", err) }
	saveQ := persist.NewSaveQueue(store, 10000)
	go func() { _ = saveQ.Run(ctx) }()

	s := zone.New(zone.Config{
		ListenAddr: listen,
		ZoneID: uint32(zoneID),
		TickHz: 20,
		AOIRadius: 25,
		CellSize: 8,
		BudgetBytes: 900,
		StateEveryTicks: 5,
		SaveEveryTicks: 20,
		Store: store,
		SaveQ: saveQ,

		TransferTargetZone: uint32(targetZone),
		TransferBoundaryX: int16(boundary),
	})
	if err := s.Start(ctx); err != nil { log.Fatalf("zone: %v", err) }
}
