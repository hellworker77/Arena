package udp_types

import (
	"sync"
	"time"
)

const (
	ackWindow = 64
)

type ClientState struct {
	LastSeqSent     uint32
	LastSeqReceived uint32
	AckBits         uint64

	LastHeard time.Time
	mu        sync.Mutex
}

func NewClientState() *ClientState {
	return &ClientState{
		LastHeard: time.Now(),
	}
}

func (c *ClientState) UpdateOnReceive(seq uint32) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.LastSeqReceived == 0 {
		c.LastSeqReceived = seq
		c.AckBits = 1
		c.LastHeard = time.Now()
		return
	}

	if seq > c.LastSeqReceived {
		shift := seq - c.LastSeqReceived

		if shift >= ackWindow {
			c.AckBits = 1
		} else {
			c.AckBits = (c.AckBits << shift) | 1
		}
	} else {
		diff := c.LastSeqReceived - seq
		if diff < ackWindow {
			c.AckBits |= 1 << diff
		}
	}

	c.LastHeard = time.Now()
}

func (c *ClientState) PrepareHeaderOnSend() (seq, ackLatest uint32, ackBitmap uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.LastSeqSent++
	return c.LastSeqSent, c.LastSeqReceived, c.AckBits
}
