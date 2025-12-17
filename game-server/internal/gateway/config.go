package gateway

import "time"

type Config struct {
	UDPListenAddr string
	Zones         ZoneFlags
	IdleTimeout   time.Duration
}
