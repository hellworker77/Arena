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
	var aoi int
	var cell int
	var budget int
	var stateEvery int
	var saveEvery int
	var storeDir string

	flag.StringVar(&listen, "listen", "127.0.0.1:4000", "TCP listen address for gateway link")
	flag.UintVar(&zoneID, "zone", 1, "Zone ID")
	flag.IntVar(&aoi, "aoi", 25, "AOI radius")
	flag.IntVar(&cell, "cell", 8, "Grid cell size")
	flag.IntVar(&budget, "budget", 900, "per-session replicate budget bytes per tick")
	flag.IntVar(&stateEvery, "stateEvery", 5, "state channel every N ticks")
	flag.IntVar(&saveEvery, "saveEvery", 20, "enqueue dirty saves every N ticks")
	flag.StringVar(&storeDir, "store", "./data", "store directory (JSON store placeholder)")
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	store, err := persist.NewJSONStore(storeDir)
	if err != nil {
		log.Fatalf("store: %v", err)
	}
	saveQ := persist.NewSaveQueue(store, 10000)
	go func() { _ = saveQ.Run(ctx) }()

	s := zone.New(zone.Config{
		ListenAddr:       listen,
		ZoneID:           uint32(zoneID),
		TickHz:           20,
		AOIRadius:        int16(aoi),
		CellSize:         int16(cell),
		MaxMoveEvents:    256,
		MaxStateEvents:   64,
		MaxEventEvents:   64,
		BudgetBytes:      budget,
		StateEveryTicks:  stateEvery,
		SaveEveryTicks:   saveEvery,
		Store:            store,
		SaveQ:            saveQ,
	})
	if err := s.Start(ctx); err != nil {
		log.Fatalf("zone: %v", err)
	}
}
