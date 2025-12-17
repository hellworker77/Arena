package gateway

import "time"

type Config struct {
	UDPListenAddr string
	ZoneTCPAddr   string
	IdleTimeout   time.Duration
}
