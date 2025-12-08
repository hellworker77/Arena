package auth

type JwtCfg struct {
	Authority              string
	Audience               string
	RotationIntervalSec    int
	AutoRefreshIntervalSec int
}
