package main

import (
	"fmt"
	"game-server/internal/udp"
	"game-server/pkg/auth"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load()
}

func main() {
	host := os.Getenv("UDP_MMO_SERVER_HOST")
	port := os.Getenv("UDP_MMO_SERVER_PORT")

	rotationIntervalSec, err := strconv.Atoi(os.Getenv("JWT_ROTATION_INTERVAL_SEC"))
	if err != nil {
		log.Fatalf("Invalid JWT_ROTATION_INTERVAL_SEC: %v", err)
	}

	autoRefreshIntervalSec, err := strconv.Atoi(os.Getenv("JWT_AUTO_REFRESH_INTERVAL_SEC"))
	if err != nil {
		log.Fatalf("Invalid JWT_AUTO_REFRESH_INTERVAL_SEC: %v", err)
	}

	jwtCfg := auth.JwtCfg{
		Authority:              os.Getenv("JWT_AUTHORITY"),
		Audience:               os.Getenv("JWT_AUTHORITY"),
		RotationIntervalSec:    rotationIntervalSec,
		AutoRefreshIntervalSec: autoRefreshIntervalSec,
	}

	if host == "" || port == "" {
		log.Fatal("UDP_MMO_SERVER_HOST or UDP_MMO_SERVER_PORT environment variables are not set in .env file")
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("Starting UDP server on %s...", addr)

	allowLegacyAuth := false
	if v := os.Getenv("ALLOW_LEGACY_AUTH"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			log.Fatalf("Invalid ALLOW_LEGACY_AUTH: %v", err)
		}
		allowLegacyAuth = b
	}

	srv, err := udp.NewServer(addr, jwtCfg, allowLegacyAuth)
	if err != nil {
		log.Fatalf("Failed to start UDP server: %v", err)
	}
	log.Printf("Listen on %s", addr)
	defer srv.Close()

	go srv.Listen()
	srv.Startup()

	log.Println("Server stopped.")
}
