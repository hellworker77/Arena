package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"game-server/internal/gateway"
)

func main() {
	var addr string
	var proto uint
	var charID uint64
	flag.StringVar(&addr, "addr", "127.0.0.1:7777", "gateway addr")
	flag.UintVar(&proto, "proto", 1, "protocol version")
	flag.Uint64Var(&charID, "char", 1, "character id")
	flag.Parse()

	raddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil { panic(err) }
	c, err := net.DialUDP("udp", nil, raddr)
	if err != nil { panic(err) }
	defer c.Close()

	peer := clientPeer{proto: uint16(proto), c: c}
	go peer.readLoop()

	// send reliable HELLO
	hello := make([]byte, 12)
	putU64(hello[0:8], charID)
	putU32(hello[8:12], uint32(0)) // interest default
	peer.sendReliable(gateway.PHello, hello)

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
			peer.sendUnreliable(gateway.PInput, pl)
			tick++
		case "a":
			if len(parts) != 3 { fmt.Println("usage: a skill targetEID"); continue }
			skill, _ := strconv.Atoi(parts[1])
			target, _ := strconv.Atoi(parts[2])
			pl := make([]byte, 10)
			putU32(pl[0:4], tick)
			putU16(pl[4:6], uint16(skill))
			putU32(pl[6:10], uint32(target))
			peer.sendReliable(gateway.PAction, pl)
			tick++
		default:
			fmt.Println("unknown")
		}
	}
}

type clientPeer struct {
	proto uint16
	c *net.UDPConn

	peer gatewayPeer
}

func (p *clientPeer) ensure() {
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

func (p *clientPeer) sendUnreliable(ptype uint8, payload []byte) {
	p.ensure()
	pkt := gateway.EncodePacket(gateway.Packet{
		Proto: p.proto, Chan: gateway.ChanUnreliable, PType: ptype,
		Seq: 0, Ack: p.peer.recvMax, AckBits: p.peer.recvMask,
		Payload: payload,
	}, nil)
	_, _ = p.c.Write(pkt)
}

func (p *clientPeer) sendReliable(ptype uint8, payload []byte) {
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

func (p *clientPeer) readLoop() {
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
			fmt.Println(string(pk.Payload))
		}
	}
}

func (p *clientPeer) retxLoop() {
	t := time.NewTicker(50 * time.Millisecond)
	defer t.Stop()
	for range t.C {
		now := time.Now()
		for seq, sm := range p.peer.pending {
			if now.Sub(sm.sentAt) >= p.peer.rto {
				if sm.retries >= 8 {
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
