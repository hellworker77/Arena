package main

import (
	"game-server/internal/config"
	"game-server/internal/udp"
	"game-server/pkg/auth"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func init() { _ = godotenv.Load() }

func main() {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	jwtCfg := auth.JwtCfg{
		Authority:              cfg.JWTAuthority,
		Audience:               cfg.JWTAudience,
		RotationIntervalSec:    cfg.JWTRotationIntervalSec,
		AutoRefreshIntervalSec: cfg.JWTAutoRefreshSec,
	}

	addr := cfg.UDPAddr()
	log.Printf("Starting UDP server on %s...", addr)

	// Session/cleanup tuning
	cliIdle := time.Duration(cfg.ClientIdleTimeoutSec) * time.Second
	cleanup := time.Duration(cfg.CleanupIntervalSec) * time.Second
	unauthTTL := time.Duration(cfg.UnauthEntryTTLSec) * time.Second

	srv, err := udp.NewServer(addr, jwtCfg, cfg.AllowLegacyAuth, cliIdle, cleanup, unauthTTL)
	if err != nil {
		log.Fatalf("Failed to start UDP server: %v", err)
	}
	defer srv.Close()

	// Replication tuning
	srv.ConfigureReplication(cfg.InterestRadius, uint32(cfg.FullSnapshotEveryTicks))

	go srv.Listen()

	// Graceful shutdown on SIGINT/SIGTERM.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigCh
		log.Printf("Shutting down...")
		srv.Close()
	}()

	srv.Startup()
}
