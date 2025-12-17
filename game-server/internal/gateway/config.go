package gateway

import "time"

type Config struct{ UDPListenAddr, ZoneTCPAddr string; IdleTimeout time.Duration }
