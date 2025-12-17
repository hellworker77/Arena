package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"game-server/internal/gateway"
)

func main() {
	var addr string
	var proto uint
	var charID uint64
	var interpMs int
	var tickHz int
	flag.StringVar(&addr, "addr", "127.0.0.1:7777", "gateway addr")
	flag.UintVar(&proto, "proto", 1, "protocol version")
	flag.Uint64Var(&charID, "char", 1, "character id")
	flag.IntVar(&interpMs, "interpMs", 150, "interpolation delay in ms")
	flag.IntVar(&tickHz, "tickHz", 20, "server tick rate (Hz)")
	flag.Parse()

	raddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil { panic(err) }
	c, err := net.DialUDP("udp", nil, raddr)
	if err != nil { panic(err) }
	defer c.Close()

	state := newClientState(uint16(proto), c, time.Duration(interpMs)*time.Millisecond, tickHz)
	go state.readLoop()
	go state.renderLoop()

	// send reliable HELLO
	hello := make([]byte, 12)
	putU64(hello[0:8], charID)
	putU32(hello[8:12], uint32(0)) // interest default
	state.sendReliable(gateway.PHello, hello)

	fmt.Println("client ready. commands:")
	fmt.Println("  m dx dy   (movement, unreliable)")
	fmt.Println("  a skill targetEID  (action, reliable)")
	fmt.Println("  q")

	in := bufio.NewScanner(os.Stdin)
	var tick uint32 = 1
	for {
		fmt.Print("> ")
		if !in.Scan() { return }
		line := strings.TrimSpace(in.Text())
		if line == "" { continue }
		if line == "q" { return }
		parts := strings.Fields(line)
		switch parts[0] {
		case "m":
			if len(parts) != 3 { fmt.Println("usage: m dx dy"); continue }
			dx, _ := strconv.Atoi(parts[1])
			dy, _ := strconv.Atoi(parts[2])
			pl := make([]byte, 8)
			putU32(pl[0:4], tick)
			putU16(pl[4:6], uint16(int16(dx)))
			putU16(pl[6:8], uint16(int16(dy)))
			state.sendUnreliable(gateway.PInput, pl)
			tick++
		case "a":
			if len(parts) != 3 { fmt.Println("usage: a skill targetEID"); continue }
			skill, _ := strconv.Atoi(parts[1])
			target, _ := strconv.Atoi(parts[2])
			pl := make([]byte, 10)
			putU32(pl[0:4], tick)
			putU16(pl[4:6], uint16(skill))
			putU32(pl[6:10], uint32(target))
			state.sendReliable(gateway.PAction, pl)
			tick++
		default:
			fmt.Println("unknown")
		}
	}
}

type sample struct {
	tick uint32
	x, y int16
	at   time.Time
}

type entityBuf struct {
	// keep last N samples sorted by tick (append-only, occasional trim)
	s []sample
}

func (b *entityBuf) add(sm sample) {
	// ignore out-of-order far behind
	if len(b.s) > 0 && sm.tick+200 < b.s[len(b.s)-1].tick {
		return
	}
	// keep monotonic: if same tick, overwrite
	if len(b.s) > 0 && b.s[len(b.s)-1].tick == sm.tick {
		b.s[len(b.s)-1] = sm
		return
	}
	b.s = append(b.s, sm)
	// trim
	if len(b.s) > 64 {
		b.s = b.s[len(b.s)-64:]
	}
}

func (b *entityBuf) interp(targetTick uint32) (x, y float64, ok bool) {
	if len(b.s) == 0 { return 0,0,false }
	// if before first
	if targetTick <= b.s[0].tick {
		return float64(b.s[0].x), float64(b.s[0].y), true
	}
	// if after last
	last := b.s[len(b.s)-1]
	if targetTick >= last.tick {
		return float64(last.x), float64(last.y), true
	}
	// find bracketing (linear scan is fine for tiny buffers)
	for i := 1; i < len(b.s); i++ {
		a := b.s[i-1]
		c := b.s[i]
		if targetTick >= a.tick && targetTick <= c.tick {
			if c.tick == a.tick {
				return float64(c.x), float64(c.y), true
			}
			t := float64(targetTick-a.tick) / float64(c.tick-a.tick)
			x = float64(a.x) + (float64(c.x)-float64(a.x))*t
			y = float64(a.y) + (float64(c.y)-float64(a.y))*t
			return x, y, true
		}
	}
	return 0,0,false
}

type clientState struct {
	proto uint16
	c *net.UDPConn

	peer gatewayPeer

	mu sync.Mutex
	interpDelay time.Duration
	tickHz int

	// clock sync (very simple): last server tick received + local time
	lastServerTick uint32
	lastServerAt time.Time

	ents map[uint32]*entityBuf
}

func newClientState(proto uint16, c *net.UDPConn, interp time.Duration, tickHz int) *clientState {
	st := &clientState{
		proto: proto,
		c: c,
		interpDelay: interp,
		tickHz: tickHz,
		ents: make(map[uint32]*entityBuf),
	}
	st.ensure()
	return st
}

func (p *clientState) ensure() {
	if p.peer.init { return }
	p.peer = gatewayPeer{
		init: true,
		nextSeq: 1,
		pending: make(map[uint32]sent),
		pendingBytes: 0,
		maxPendingBytes: 65536,
		recvMax: 0,
		recvMask: 0,
		rto: 200*time.Millisecond,
	}
	go p.retxLoop()
}

func (p *clientState) sendUnreliable(ptype uint8, payload []byte) {
	p.ensure()
	pkt := gateway.EncodePacket(gateway.Packet{
		Proto: p.proto, Chan: gateway.ChanUnreliable, PType: ptype,
		Seq: 0, Ack: p.peer.recvMax, AckBits: p.peer.recvMask,
		Payload: payload,
	}, nil)
	_, _ = p.c.Write(pkt)
}

func (p *clientState) sendReliable(ptype uint8, payload []byte) {
	p.ensure()
	seq := p.peer.nextSeq
	p.peer.nextSeq++
	pkt := gateway.EncodePacket(gateway.Packet{
		Proto: p.proto, Chan: gateway.ChanReliable, PType: ptype,
		Seq: seq, Ack: p.peer.recvMax, AckBits: p.peer.recvMask,
		Payload: payload,
	}, nil)
	if p.peer.pendingBytes+len(pkt) > p.peer.maxPendingBytes {
		fmt.Println("reliable backlog overflow (client)")
		return
	}
	p.peer.pending[seq] = sent{pkt: pkt, sentAt: time.Now(), retries: 0}
	p.peer.pendingBytes += len(pkt)
	_, _ = p.c.Write(pkt)
}

func (p *clientState) readLoop() {
	buf := make([]byte, 64*1024)
	for {
		n, err := p.c.Read(buf)
		if err != nil { return }
		p.ensure()
		pk, err := gateway.DecodePacket(buf[:n])
		if err != nil { continue }
		if pk.Proto != p.proto { continue }

		// ack processing
		p.peer.onAcks(pk.Ack, pk.AckBits)
		if pk.Chan == gateway.ChanReliable {
			p.peer.updateRecv(pk.Seq)
		}

		switch pk.PType {
		case gateway.PText:
			fmt.Println(string(pk.Payload))
		case gateway.PRep:
			line := string(pk.Payload)
			fmt.Println(line)
			p.consumeRepLine(line)
		}
	}
}

// Rep lines are prefixed by gateway as:
// "T <serverTick> <rest...>"
func (p *clientState) consumeRepLine(line string) {
	parts := strings.Fields(line)
	if len(parts) < 3 || parts[0] != "T" {
		return
	}
	tick64, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil { return }
	serverTick := uint32(tick64)

	p.mu.Lock()
	p.lastServerTick = serverTick
	p.lastServerAt = time.Now()
	p.mu.Unlock()

	kind := parts[2]
	switch kind {
	case "SPAWN":
		if len(parts) < 6 { return }
		eid, _ := strconv.ParseUint(parts[3], 10, 32)
		x64, _ := strconv.ParseInt(parts[4], 10, 16)
		y64, _ := strconv.ParseInt(parts[5], 10, 16)
		p.addSample(uint32(eid), serverTick, int16(x64), int16(y64))
	case "MOV":
		if len(parts) < 6 { return }
		eid, _ := strconv.ParseUint(parts[3], 10, 32)
		x64, _ := strconv.ParseInt(parts[4], 10, 16)
		y64, _ := strconv.ParseInt(parts[5], 10, 16)
		p.addSample(uint32(eid), serverTick, int16(x64), int16(y64))
	case "DESPAWN":
		if len(parts) < 4 { return }
		eid, _ := strconv.ParseUint(parts[3], 10, 32)
		p.mu.Lock()
		delete(p.ents, uint32(eid))
		p.mu.Unlock()
	}
}

func (p *clientState) addSample(eid uint32, tick uint32, x, y int16) {
	p.mu.Lock()
	b := p.ents[eid]
	if b == nil {
		b = &entityBuf{}
		p.ents[eid] = b
	}
	b.add(sample{tick: tick, x: x, y: y, at: time.Now()})
	p.mu.Unlock()
}

func (p *clientState) renderLoop() {
	// Render at 10Hz just to demonstrate smoothing.
	t := time.NewTicker(100 * time.Millisecond)
	defer t.Stop()

	tickDur := time.Second
	if p.tickHz > 0 {
		tickDur = time.Second / time.Duration(p.tickHz)
	}
	for range t.C {
		p.mu.Lock()
		lastTick := p.lastServerTick
		lastAt := p.lastServerAt
		delay := p.interpDelay
		// if no sync yet
		if lastTick == 0 || lastAt.IsZero() {
			p.mu.Unlock()
			continue
		}
		// estimate current server tick based on elapsed local time since last packet
		el := time.Since(lastAt)
		estNowTick := lastTick + uint32(el / tickDur)
		// render slightly behind
		renderTick := estNowTick
		if delay > 0 {
			behind := uint32(delay / tickDur)
			if renderTick > behind {
				renderTick -= behind
			}
		}

		for eid, b := range p.ents {
			x, y, ok := b.interp(renderTick)
			if !ok { continue }
			fmt.Printf("RENDER tick=%d eid=%d x=%.2f y=%.2f
", renderTick, eid, x, y)
		}
		p.mu.Unlock()
	}
}

func (p *clientState) retxLoop() {
	t := time.NewTicker(50 * time.Millisecond)
	defer t.Stop()
	for range t.C {
		now := time.Now()
		for seq, sm := range p.peer.pending {
			if now.Sub(sm.sentAt) >= p.peer.rto {
				if sm.retries >= 8 {
					p.peer.pendingBytes -= len(p.peer.pending[seq].pkt)
					if p.peer.pendingBytes < 0 { p.peer.pendingBytes = 0 }
					delete(p.peer.pending, seq)
					continue
				}
				sm.retries++
				sm.sentAt = now
				p.peer.pending[seq] = sm
				_, _ = p.c.Write(sm.pkt)
			}
		}
	}
}

type sent struct {
	pkt []byte
	sentAt time.Time
	retries int
}

type gatewayPeer struct {
	init bool
	nextSeq uint32
	pending map[uint32]sent
	pendingBytes int
	maxPendingBytes int

	recvMax uint32
	recvMask uint32

	rto time.Duration
}

func (p *gatewayPeer) updateRecv(seq uint32) {
	if seq == 0 { return }
	if p.recvMax == 0 {
		p.recvMax = seq; p.recvMask = 0; return
	}
	if seq > p.recvMax {
		shift := seq - p.recvMax
		if shift >= 32 { p.recvMask = 0 } else { p.recvMask = (p.recvMask << shift) | (1 << (shift-1)) }
		p.recvMax = seq
		return
	}
	d := p.recvMax - seq
	if d == 0 { return }
	if d <= 32 { p.recvMask |= 1 << (d-1) }
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

func (p *gatewayPeer) onAcks(ack uint32, ackBits uint32) {
	for seq := range p.pending {
		if ackedBy(ack, ackBits, seq) {
			p.pendingBytes -= len(p.pending[seq].pkt)
			if p.pendingBytes < 0 { p.pendingBytes = 0 }
			delete(p.pending, seq)
		}
	}
}

func putU16(b []byte, v uint16) { b[0]=byte(v); b[1]=byte(v>>8) }
func putU32(b []byte, v uint32) { putU16(b[0:2], uint16(v)); putU16(b[2:4], uint16(v>>16)) }
func putU64(b []byte, v uint64) { putU32(b[0:4], uint32(v)); putU32(b[4:8], uint32(v>>32)) }
