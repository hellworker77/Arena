package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"math/big"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

type JwksCache struct {
	mu                  sync.RWMutex
	keys                map[string]*jwt.SigningMethodRSA
	decodingKeys        map[string]*rsa.PublicKey
	authority           string
	rotationInterval    time.Duration
	autoRefreshInterval time.Duration
}

func NewJwksCache(authority string, rotateSec, refreshSec int) (*JwksCache, error) {
	cache := &JwksCache{
		keys:                make(map[string]*jwt.SigningMethodRSA),
		decodingKeys:        make(map[string]*rsa.PublicKey),
		authority:           authority,
		rotationInterval:    time.Duration(rotateSec) * time.Second,
		autoRefreshInterval: time.Duration(refreshSec) * time.Second,
	}

	if err := cache.Refresh(); err != nil {
		return nil, err
	}

	return cache, nil
}

func rsaFromComponents(nStr, eStr string) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(nStr)
	if err != nil {
		return nil, err
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(eStr)
	if err != nil {
		return nil, err
	}

	n := new(big.Int).SetBytes(nBytes)

	var e big.Int
	e.SetBytes(eBytes)

	pub := &rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}

	return pub, nil
}

func (cache *JwksCache) Refresh() error {
	jwks, err := FetchJwks(cache.authority)

	if err != nil {
		log.Warn().Err(err).Msg("failed to refresh JWKS")
		return err
	}

	cache.mu.Lock()
	defer cache.mu.Unlock()

	for _, key := range jwks.Keys {
		publicKey, err := rsaFromComponents(key.N, key.E)
		if err != nil {
			log.Warn().Err(err).Msgf("invalid JWKS key %s", key.Kid)
			continue
		}

		cache.decodingKeys[key.Kid] = publicKey
	}

	log.Info().Msg("JWKS refreshed successfully")
	return nil
}

func (cache *JwksCache) autoRefresh() {
	ticker := time.NewTicker(cache.autoRefreshInterval)
	defer ticker.Stop()

	for range ticker.C {
		cache.Refresh()
	}
}

func (cache *JwksCache) GetKey(kid string) (*rsa.PublicKey, bool) {
	cache.mu.RLock()
	defer cache.mu.RUnlock()
	k, ok := cache.decodingKeys[kid]
	return k, ok
}
