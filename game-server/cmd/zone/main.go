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
	var httpAddr string
	var zoneID uint
	var storeDir string

	flag.StringVar(&listen, "listen", "127.0.0.1:4000", "TCP listen address for gateway link")
	flag.StringVar(&httpAddr, "http", "", "HTTP metrics address (e.g. :9101)")
	flag.UintVar(&zoneID, "zone", 1, "Zone ID")
	flag.StringVar(&storeDir, "store", "./data", "store directory")
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	store, err := persist.NewJSONStore(storeDir)
	if err != nil { log.Fatalf("store: %v", err) }
	saveQ := persist.NewSaveQueue(store, 10000)
	go func() { _ = saveQ.Run(ctx) }()

	snapStore, err := persist.NewJSONSnapshotStore(storeDir)
	if err != nil { log.Fatalf("snapshot store: %v", err) }
	snapQ := persist.NewSnapshotQueue(snapStore, 1000)
	go func() { _ = snapQ.Run(ctx) }()

	// toy transfer mapping:
	// zone 1 transfers to 2 when X > 100
	// zone 2 transfers to 1 when X < -100
	var target uint32 = 2
	var boundary int16 = 100
	if zoneID == 2 {
		target = 1
		boundary = -100
	}

	s := zone.New(zone.Config{
		ListenAddr: listen,
		HTTPAddr: httpAddr,
		ZoneID: uint32(zoneID),
		TickHz: 20,
		AOIRadius: 25,
		CellSize: 8,
		BudgetBytes: 900,
		StateEveryTicks: 5,
		SaveEveryTicks: 20,
		Store: store,
		SaveQ: saveQ,
		SnapshotEveryTicks: 200,
		SnapshotStore: snapStore,
		SnapshotQ: snapQ,
		AIBudgetPerTick: 200,
		TransferTargetZone: target,
		TransferBoundaryX: boundary,
		TransferTimeoutTicks: 60,
		HistoryTicks: 40,
		RewindMaxTicks: 5,
	})
	if err := s.Start(ctx); err != nil { log.Fatalf("zone: %v", err) }
}
