package gateway

import "time"

type Config struct {
	UDPListenAddr string
	Zones         ZoneFlags
	IdleTimeout   time.Duration
	ProtoVersion  uint16
	TransferTimeout time.Duration

	RateBytesPerSec int
	BurstBytes       int
	MaxReliableBytes int
}
