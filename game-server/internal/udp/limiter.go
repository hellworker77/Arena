package udp

import "time"

// tokenBucket is a tiny per-client rate limiter.
// Not thread-safe; use from one goroutine (packet processor).
type tokenBucket struct {
	capacity  float64
	tokens    float64
	refillPer float64 // tokens per second
	last      time.Time
}

func newTokenBucket(capacity int, refillPerSec int) *tokenBucket {
	capF := float64(capacity)
	return &tokenBucket{
		capacity:  capF,
		tokens:    capF,
		refillPer: float64(refillPerSec),
		last:      time.Now(),
	}
}

// allow consumes n tokens if possible and returns true. Otherwise returns false.
func (b *tokenBucket) allow(n float64) bool {
	now := time.Now()
	dt := now.Sub(b.last).Seconds()
	if dt > 0 {
		b.tokens += dt * b.refillPer
		if b.tokens > b.capacity {
			b.tokens = b.capacity
		}
		b.last = now
	}
	if b.tokens < n {
		return false
	}
	b.tokens -= n
	return true
}
