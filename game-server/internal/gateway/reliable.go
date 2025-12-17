package gateway

import (
	"time"
)

type sentMsg struct {
	seq uint32
	pkt []byte
	sentAt time.Time
	retries int
}

type reliablePeer struct {
	// send
	nextSeq uint32
	pending map[uint32]*sentMsg

	// recv
	recvMax uint32
	recvMask uint32 // bits for last 32 packets behind recvMax

	rto time.Duration
	maxRetries int
}

func newPeer() *reliablePeer {
	return &reliablePeer{
		nextSeq: 1,
		pending: make(map[uint32]*sentMsg),
		recvMax: 0,
		recvMask: 0,
		rto: 200 * time.Millisecond,
		maxRetries: 8,
	}
}

// updateRecv tracks which reliable seq we've seen.
func (p *reliablePeer) updateRecv(seq uint32) {
	if seq == 0 {
		return
	}
	if p.recvMax == 0 {
		p.recvMax = seq
		p.recvMask = 0
		return
	}
	if seq > p.recvMax {
		shift := seq - p.recvMax
		if shift >= 32 {
			p.recvMask = 0
		} else {
			p.recvMask = (p.recvMask << shift) | (1 << (shift-1))
		}
		p.recvMax = seq
		return
	}
	// seq <= recvMax
	d := p.recvMax - seq
	if d == 0 {
		return
	}
	if d <= 32 {
		p.recvMask |= 1 << (d-1)
	}
}

func ackedBy(ack uint32, ackBits uint32, seq uint32) bool {
	if seq == 0 { return true }
	if ack == seq { return true }
	if seq > ack { return false }
	d := ack - seq
	if d == 0 { return true }
	if d > 32 { return false }
	return (ackBits & (1 << (d-1))) != 0
}

func (p *reliablePeer) onAcks(ack uint32, ackBits uint32) {
	for seq := range p.pending {
		if ackedBy(ack, ackBits, seq) {
			delete(p.pending, seq)
		}
	}
}

func (p *reliablePeer) allocSeq() uint32 {
	s := p.nextSeq
	p.nextSeq++
	if p.nextSeq == 0 { p.nextSeq = 1 }
	return s
}
