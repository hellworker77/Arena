package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	UDPHost string
	UDPPort string

	JWTAuthority           string
	JWTAudience            string
	JWTRotationIntervalSec int
	JWTAutoRefreshSec      int

	// If true, accepts legacy PTAuth payloads (raw JWT string without handshake cookie).
	// Keep this false in production.
	AllowLegacyAuth bool

	// Session/cleanup settings
	ClientIdleTimeoutSec int
	CleanupIntervalSec   int
	UnauthEntryTTLSec    int
}

func LoadFromEnv() (Config, error) {
	var c Config
	c.UDPHost = os.Getenv("UDP_MMO_SERVER_HOST")
	c.UDPPort = os.Getenv("UDP_MMO_SERVER_PORT")
	c.JWTAuthority = os.Getenv("JWT_AUTHORITY")
	c.JWTAudience = os.Getenv("JWT_AUDIENCE")
	if c.JWTAudience == "" {
		// Backward-compat with previous .env where Audience mirrored Authority.
		c.JWTAudience = os.Getenv("JWT_AUTHORITY")
	}

	rot, err := strconv.Atoi(os.Getenv("JWT_ROTATION_INTERVAL_SEC"))
	if err != nil {
		return c, fmt.Errorf("invalid JWT_ROTATION_INTERVAL_SEC: %w", err)
	}
	ref, err := strconv.Atoi(os.Getenv("JWT_AUTO_REFRESH_INTERVAL_SEC"))
	if err != nil {
		return c, fmt.Errorf("invalid JWT_AUTO_REFRESH_INTERVAL_SEC: %w", err)
	}
	c.JWTRotationIntervalSec = rot
	c.JWTAutoRefreshSec = ref

	// Optional toggles
	if v := os.Getenv("ALLOW_LEGACY_AUTH"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return c, fmt.Errorf("invalid ALLOW_LEGACY_AUTH: %w", err)
		}
		c.AllowLegacyAuth = b
	}

	// Optional session/cleanup tuning (safe defaults)
	// CLIENT_IDLE_TIMEOUT_SEC: authenticated client is disconnected if no packets are heard within this window.
	// CLEANUP_INTERVAL_SEC: how often server prunes idle sessions and stale limiter state.
	// UNAUTH_ENTRY_TTL_SEC: how long to keep pre-auth limiter/auth-attempt state for an address.
	c.ClientIdleTimeoutSec = 30
	c.CleanupIntervalSec = 5
	c.UnauthEntryTTLSec = 60
	if v := os.Getenv("CLIENT_IDLE_TIMEOUT_SEC"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return c, fmt.Errorf("invalid CLIENT_IDLE_TIMEOUT_SEC: %w", err)
		}
		c.ClientIdleTimeoutSec = n
	}
	if v := os.Getenv("CLEANUP_INTERVAL_SEC"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return c, fmt.Errorf("invalid CLEANUP_INTERVAL_SEC: %w", err)
		}
		c.CleanupIntervalSec = n
	}
	if v := os.Getenv("UNAUTH_ENTRY_TTL_SEC"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return c, fmt.Errorf("invalid UNAUTH_ENTRY_TTL_SEC: %w", err)
		}
		c.UnauthEntryTTLSec = n
	}

	if c.UDPHost == "" || c.UDPPort == "" {
		return c, fmt.Errorf("UDP_MMO_SERVER_HOST/UDP_MMO_SERVER_PORT must be set")
	}
	if c.JWTAuthority == "" {
		return c, fmt.Errorf("JWT_AUTHORITY must be set")
	}
	return c, nil
}

func (c Config) UDPAddr() string { return fmt.Sprintf("%s:%s", c.UDPHost, c.UDPPort) }
