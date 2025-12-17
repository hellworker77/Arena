package gateway

import "time"

type sentMsg struct {
	seq     uint32
	pkt     []byte
	sentAt  time.Time
	retries int
	size    int
}

type rttEstimator struct {
	srtt   time.Duration
	rttvar time.Duration
	rto    time.Duration
	inited bool
}

func (r *rttEstimator) update(sample time.Duration) {
	if sample <= 0 {
		return
	}
	// Jacobson/Karels (RFC6298-ish constants)
	const (
		alpha = 1.0 / 8.0
		beta  = 1.0 / 4.0
	)
	if !r.inited {
		r.srtt = sample
		r.rttvar = sample / 2
		r.rto = clamp(sample*3, 200*time.Millisecond, 2*time.Second)
		r.inited = true
		return
	}
	// rttvar = (1-beta)*rttvar + beta*|srtt-sample|
	if r.srtt > sample {
		r.rttvar = time.Duration((1-beta)*float64(r.rttvar) + beta*float64(r.srtt-sample))
	} else {
		r.rttvar = time.Duration((1-beta)*float64(r.rttvar) + beta*float64(sample-r.srtt))
	}
	// srtt = (1-alpha)*srtt + alpha*sample
	r.srtt = time.Duration((1-alpha)*float64(r.srtt) + alpha*float64(sample))
	r.rto = clamp(r.srtt+4*r.rttvar, 200*time.Millisecond, 2*time.Second)
}

func clamp(v, lo, hi time.Duration) time.Duration {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

type reliablePeer struct {
	// send
	nextSeq      uint32
	pending      map[uint32]*sentMsg
	pendingBytes int
	maxPendingBytes int

	// recv
	recvMax  uint32
	recvMask uint32 // bits for last 32 packets behind recvMax

	est rttEstimator
	maxRetries int
}

func newPeer(maxPendingBytes int) *reliablePeer {
	if maxPendingBytes <= 0 {
		maxPendingBytes = 65536
	}
	p := &reliablePeer{
		nextSeq:  1,
		pending:  make(map[uint32]*sentMsg),
		recvMax:  0,
		recvMask: 0,
		maxRetries: 8,
		maxPendingBytes: maxPendingBytes,
	}
	// default RTO before first sample
	p.est.rto = 200 * time.Millisecond
	return p
}

func (p *reliablePeer) currentRTO() time.Duration {
	if p.est.rto <= 0 {
		return 200 * time.Millisecond
	}
	return p.est.rto
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
			p.recvMask = (p.recvMask << shift) | (1 << (shift - 1))
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
		p.recvMask |= 1 << (d - 1)
	}
}

func ackedBy(ack uint32, ackBits uint32, seq uint32) bool {
	if seq == 0 {
		return true
	}
	if ack == seq {
		return true
	}
	if seq > ack {
		return false
	}
	d := ack - seq
	if d == 0 {
		return true
	}
	if d > 32 {
		return false
	}
	return (ackBits & (1 << (d - 1))) != 0
}

func (p *reliablePeer) onAcks(now time.Time, ack uint32, ackBits uint32) {
	for seq, sm := range p.pending {
		if ackedBy(ack, ackBits, seq) {
			// RTT sample from first-send time of this msg
			if !sm.sentAt.IsZero() && sm.retries == 0 {
				p.est.update(now.Sub(sm.sentAt))
			}
			p.pendingBytes -= sm.size
			if p.pendingBytes < 0 {
				p.pendingBytes = 0
			}
			delete(p.pending, seq)
		}
	}
}

func (p *reliablePeer) allocSeq() uint32 {
	s := p.nextSeq
	p.nextSeq++
	if p.nextSeq == 0 {
		p.nextSeq = 1
	}
	return s
}

func (p *reliablePeer) canEnqueue(size int) bool {
	return p.pendingBytes+size <= p.maxPendingBytes
}

func (p *reliablePeer) enqueue(seq uint32, pkt []byte) {
	size := len(pkt)
	p.pending[seq] = &sentMsg{seq: seq, pkt: pkt, sentAt: time.Now(), retries: 0, size: size}
	p.pendingBytes += size
}
