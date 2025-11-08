package main

import (
	"fmt"
	"game-server/internal/infrastructure/udp"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load()
}

func main() {
	host := os.Getenv("UDP_MMO_SERVER_HOST")
	port := os.Getenv("UDP_MMO_SERVER_PORT")

	if host == "" || port == "" {
		log.Fatal("UDP_MMO_SERVER_HOST or UDP_MMO_SERVER_PORT environment variables are not set in .env file")
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("Starting UDP server on %s...", addr)

	srv, err := udp.NewServer(addr)
	if err != nil {
		log.Fatalf("Failed to start UDP server: %v", err)
	}
	log.Printf("Listen on %s", addr)
	defer srv.Close()

	go srv.Listen()
	srv.Startup()

	log.Println("Server stopped.")
}
