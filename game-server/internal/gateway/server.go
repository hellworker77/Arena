package gateway

import (
	"bufio"
	"context"
	"errors"
	"log"
	"net"
	"sync"
	"time"

	"game-server/internal/shared"
	"game-server/internal/shared/wire"
)

type Server struct {
	cfg Config
	udpConn *net.UDPConn

	zonesMu sync.Mutex
	zones map[uint32]*zoneLink

	sessionsMu sync.Mutex
	byRemote map[string]*sessionState
	bySID map[shared.SessionID]string

	// transfer inflight (Step13)
	xferMu sync.Mutex
	inflight map[shared.SessionID]*xferState
}

type zoneLink struct {
	id uint32
	addr string
	conn net.Conn
	r *bufio.Reader
	w *bufio.Writer
}

type sessionState struct {
	SID shared.SessionID
	CharID shared.CharacterID
	ZoneID shared.ZoneID
	Interest wire.InterestMask
	LastHeard time.Time
	Proto uint16

	peer *reliablePeer
	raddr *net.UDPAddr
}

type xferState struct {
	From shared.ZoneID
	To shared.ZoneID
	Started time.Time
}

func New(cfg Config) (*Server, error) {
	if cfg.UDPListenAddr == "" { return nil, errors.New("missing udp addr") }
	if len(cfg.Zones) == 0 { return nil, errors.New("no zones configured") }
	if cfg.IdleTimeout <= 0 { cfg.IdleTimeout = 30*time.Second }
	if cfg.TransferTimeout <= 0 { cfg.TransferTimeout = 3*time.Second }
	if cfg.ProtoVersion == 0 { cfg.ProtoVersion = 1 }

	return &Server{
		cfg: cfg,
		zones: make(map[uint32]*zoneLink),
		byRemote: make(map[string]*sessionState),
		bySID: make(map[shared.SessionID]string),
		inflight: make(map[shared.SessionID]*xferState),
	}, nil
}

func (s *Server) Start(ctx context.Context) error {
	udpAddr, err := net.ResolveUDPAddr("udp", s.cfg.UDPListenAddr)
	if err != nil { return err }
	s.udpConn, err = net.ListenUDP("udp", udpAddr)
	if err != nil { return err }
	defer s.udpConn.Close()

	for zid, addr := range s.cfg.Zones {
		if err := s.connectZone(zid, addr); err != nil { return err }
	}
	defer s.closeZones()

	s.zonesMu.Lock()
	for _, zl := range s.zones {
		go s.zoneReadLoop(ctx, zl)
	}
	s.zonesMu.Unlock()

	go s.cleanupLoop(ctx)
	go s.transferTimeoutLoop(ctx)
	go s.retransmitLoop(ctx)

	log.Printf("gateway up: udp=%s zones=%d proto=%d", s.cfg.UDPListenAddr, len(s.cfg.Zones), s.cfg.ProtoVersion)

	buf := make([]byte, 64*1024)
	for {
		select { case <-ctx.Done(): return nil; default: }
		_ = s.udpConn.SetReadDeadline(time.Now().Add(500*time.Millisecond))
		n, raddr, err := s.udpConn.ReadFromUDP(buf)
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() { continue }
			return err
		}
		s.handleUDPPacket(raddr, buf[:n])
	}
}

func (s *Server) connectZone(zid uint32, addr string) error {
	c, err := net.Dial("tcp", addr)
	if err != nil { return err }
	zl := &zoneLink{
		id: zid, addr: addr,
		conn: c,
		r: bufio.NewReaderSize(c, 64*1024),
		w: bufio.NewWriterSize(c, 64*1024),
	}
	s.zonesMu.Lock()
	s.zones[zid] = zl
	s.zonesMu.Unlock()
	return nil
}

func (s *Server) closeZones() {
	s.zonesMu.Lock()
	defer s.zonesMu.Unlock()
	for _, zl := range s.zones { _ = zl.conn.Close() }
	s.zones = make(map[uint32]*zoneLink)
}

func (s *Server) zoneSend(zid uint32, typ wire.MsgType, payload []byte) error {
	s.zonesMu.Lock()
	zl := s.zones[zid]
	s.zonesMu.Unlock()
	if zl == nil { return errors.New("unknown zone") }
	return wire.WriteFrame(zl.w, typ, payload)
}

func (s *Server) handleUDPPacket(raddr *net.UDPAddr, b []byte) {
	p, err := DecodePacket(b)
	if err != nil {
		return
	}
	if p.Proto != s.cfg.ProtoVersion {
		return
	}
	remote := raddr.String()

	st := s.getOrCreate(remote, raddr, p.Proto)

	// Process ACKs for reliable channel (both on any packet)
	st.peer.onAcks(p.Ack, p.AckBits)

	// Update receive window if incoming is reliable
	if p.Chan == ChanReliable {
		st.peer.updateRecv(p.Seq)
	}

	st.LastHeard = time.Now()

	switch p.PType {
	case PHello:
		if p.Chan != ChanReliable {
			return
		}
		if len(p.Payload) < 8+4 { return }
		cid := shared.CharacterID(binaryLEU64(p.Payload[0:8]))
		interest := wire.InterestMask(binaryLEU32(p.Payload[8:12]))
		if cid == 0 { return }
		if interest == 0 {
			interest = wire.InterestMove | wire.InterestState | wire.InterestEvent | wire.InterestCombat
		}
		st.CharID = cid
		st.Interest = interest

		// choose default zone (min id)
		if st.ZoneID == 0 {
			var min uint32 = 0
			for zid := range s.cfg.Zones {
				if min == 0 || zid < min { min = zid }
			}
			st.ZoneID = shared.ZoneID(min)
		}
		_ = s.zoneSend(uint32(st.ZoneID), wire.MsgAttachPlayer, wire.EncodeAttachPlayer(st.SID, st.CharID, st.ZoneID, st.Interest))

		// Send a reliable text ACK to client
		s.sendReliableText(st, "HELLO_OK sid="+st.SID.String())

	case PInput:
		if len(p.Payload) < 4+2+2 { return }
		if st.ZoneID == 0 { return }
		tick := binaryLEU32(p.Payload[0:4])
		mx := int16(binaryLEU16(p.Payload[4:6]))
		my := int16(binaryLEU16(p.Payload[6:8]))
		_ = s.zoneSend(uint32(st.ZoneID), wire.MsgPlayerInput, wire.EncodePlayerInput(st.SID, tick, mx, my))

	case PAction:
		if p.Chan != ChanReliable { return }
		if len(p.Payload) < 4+2+4 { return }
		if st.ZoneID == 0 { return }
		tick := binaryLEU32(p.Payload[0:4])
		skill := binaryLEU16(p.Payload[4:6])
		target := shared.EntityID(binaryLEU32(p.Payload[6:10]))
		_ = s.zoneSend(uint32(st.ZoneID), wire.MsgPlayerAction, wire.EncodePlayerAction(st.SID, tick, skill, target))

	default:
	}
}

func (s *Server) getOrCreate(remote string, raddr *net.UDPAddr, proto uint16) *sessionState {
	s.sessionsMu.Lock()
	defer s.sessionsMu.Unlock()
	if st, ok := s.byRemote[remote]; ok {
		st.raddr = raddr
		return st
	}
	st := &sessionState{
		SID: shared.NewSessionID(),
		CharID: 0,
		ZoneID: 0,
		Interest: 0,
		LastHeard: time.Now(),
		Proto: proto,
		peer: newPeer(),
		raddr: raddr,
	}
	s.byRemote[remote] = st
	s.bySID[st.SID] = remote
	return st
}

func (s *Server) getBySID(sid shared.SessionID) (*sessionState, bool) {
	s.sessionsMu.Lock()
	defer s.sessionsMu.Unlock()
	ra, ok := s.bySID[sid]
	if !ok { return nil, false }
	st := s.byRemote[ra]
	return st, st != nil
}

func (s *Server) cleanupLoop(ctx context.Context) {
	t := time.NewTicker(2*time.Second)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			now := time.Now()
			var toDetach []struct{ sid shared.SessionID; zid shared.ZoneID }
			s.sessionsMu.Lock()
			for ra, st := range s.byRemote {
				if now.Sub(st.LastHeard) > s.cfg.IdleTimeout {
					toDetach = append(toDetach, struct{sid shared.SessionID; zid shared.ZoneID}{st.SID, st.ZoneID})
					delete(s.bySID, st.SID)
					delete(s.byRemote, ra)
				}
			}
			s.sessionsMu.Unlock()
			for _, d := range toDetach {
				if d.zid != 0 {
					_ = s.zoneSend(uint32(d.zid), wire.MsgDetachPlayer, wire.EncodeDetachPlayer(d.sid))
				}
			}
		}
	}
}

func (s *Server) transferTimeoutLoop(ctx context.Context) {
	t := time.NewTicker(250 * time.Millisecond)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			now := time.Now()
			var abort []struct{ sid shared.SessionID; from shared.ZoneID }
			s.xferMu.Lock()
			for sid, xs := range s.inflight {
				if now.Sub(xs.Started) > s.cfg.TransferTimeout {
					abort = append(abort, struct{sid shared.SessionID; from shared.ZoneID}{sid, xs.From})
					delete(s.inflight, sid)
				}
			}
			s.xferMu.Unlock()
			for _, a := range abort {
				_ = s.zoneSend(uint32(a.from), wire.MsgTransferAbort, wire.EncodeTransferAbort(a.sid))
				if st, ok := s.getBySID(a.sid); ok {
					s.sendReliableText(st, "XFER_ABORT timeout")
				}
			}
		}
	}
}

func (s *Server) retransmitLoop(ctx context.Context) {
	t := time.NewTicker(50 * time.Millisecond)
	defer t.Stop()
	buf := make([]byte, 0, 2048)
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			now := time.Now()
			// iterate sessions and resend pending
			var sends []struct{
				addr *net.UDPAddr
				pkt []byte
				st *sessionState
				seq uint32
			}
			s.sessionsMu.Lock()
			for _, st := range s.byRemote {
				if st.raddr == nil || st.peer == nil { continue }
				for seq, sm := range st.peer.pending {
					if now.Sub(sm.sentAt) >= st.peer.rto {
						if sm.retries >= st.peer.maxRetries {
							delete(st.peer.pending, seq)
							continue
						}
						sm.retries++
						sm.sentAt = now
						sends = append(sends, struct{
							addr *net.UDPAddr; pkt []byte; st *sessionState; seq uint32
						}{addr: st.raddr, pkt: sm.pkt, st: st, seq: seq})
					}
				}
			}
			s.sessionsMu.Unlock()
			for _, it := range sends {
				_, _ = s.udpConn.WriteToUDP(it.pkt, it.addr)
				_ = buf
			}
		}
	}
}

func (s *Server) sendReliableText(st *sessionState, msg string) {
	if st == nil || st.raddr == nil { return }
	seq := st.peer.allocSeq()
	payload := []byte(msg)
	pkt := EncodePacket(Packet{
		Proto: s.cfg.ProtoVersion,
		Chan: ChanReliable,
		PType: PText,
		Seq: seq,
		Ack: st.peer.recvMax,
		AckBits: st.peer.recvMask,
		Payload: payload,
	}, nil)
	st.peer.pending[seq] = &sentMsg{seq: seq, pkt: pkt, sentAt: time.Now(), retries: 0}
	_, _ = s.udpConn.WriteToUDP(pkt, st.raddr)
}

func (s *Server) sendUnreliableRep(st *sessionState, line string) {
	if st == nil || st.raddr == nil { return }
	pkt := EncodePacket(Packet{
		Proto: s.cfg.ProtoVersion,
		Chan: ChanUnreliable,
		PType: PRep,
		Seq: 0,
		Ack: st.peer.recvMax,
		AckBits: st.peer.recvMask,
		Payload: []byte(line),
	}, nil)
	_, _ = s.udpConn.WriteToUDP(pkt, st.raddr)
}

func (s *Server) zoneReadLoop(ctx context.Context, zl *zoneLink) {
	for {
		select { case <-ctx.Done(): return; default: }
		fr, err := wire.ReadFrame(zl.r)
		if err != nil {
			log.Printf("zone %d link read error: %v", zl.id, err)
			return
		}
		switch fr.Type {
		case wire.MsgAttachAck:
			s.tryCommitOnAttachAck(zl.id)
		case wire.MsgReplicate:
			sid, _, ch, events, err := wire.DecodeReplicate(fr.Payload)
			if err != nil { continue }
			st, ok := s.getBySID(sid)
			if !ok { continue }
			// ship as human-readable lines (demo), unreliable
			switch ch {
			case wire.ChanMove:
				for _, ev := range events {
					switch ev.Op {
					case wire.RepSpawn:
						s.sendUnreliableRep(st, sprintf("SPAWN %d %d %d kind=%d mask=%d", uint32(ev.EID), ev.X, ev.Y, ev.Kind, uint32(ev.Mask)))
					case wire.RepDespawn:
						s.sendUnreliableRep(st, sprintf("DESPAWN %d", uint32(ev.EID)))
					case wire.RepMove:
						s.sendUnreliableRep(st, sprintf("MOV %d %d %d", uint32(ev.EID), ev.X, ev.Y))
					}
				}
			case wire.ChanState:
				for _, ev := range events {
					if ev.Op == wire.RepStateHP {
						s.sendUnreliableRep(st, sprintf("STAT %d hp=%d", uint32(ev.EID), ev.Val))
					}
				}
			case wire.ChanEvent:
				for _, ev := range events {
					if ev.Op == wire.RepEventText {
						// demo: send reliable event text
						s.sendReliableText(st, "EV "+ev.Text)
					}
				}
			}
		case wire.MsgTransferPrepare:
			sid, _, target, interest, x, y, hp, err := wire.DecodeTransferPrepare(fr.Payload)
			if err != nil { continue }

			st, ok := s.getBySID(sid)
			if !ok { continue }

			// validate target exists
			s.zonesMu.Lock()
			_, okTarget := s.zones[uint32(target)]
			s.zonesMu.Unlock()
			if !okTarget {
				_ = s.zoneSend(uint32(st.ZoneID), wire.MsgTransferAbort, wire.EncodeTransferAbort(sid))
				s.sendReliableText(st, "XFER_ABORT bad_target")
				continue
			}

			s.xferMu.Lock()
			s.inflight[sid] = &xferState{
				From: st.ZoneID,
				To: shared.ZoneID(target),
				Started: time.Now(),
			}
			s.xferMu.Unlock()

			// attach to target
			_ = s.zoneSend(uint32(target), wire.MsgAttachWithState, wire.EncodeAttachWithState(sid, st.CharID, shared.ZoneID(target), interest, x, y, hp))
			// route to target
			st.ZoneID = shared.ZoneID(target)
			st.Interest = interest
			s.sendReliableText(st, sprintf("XFER_PREP %d->%d", zl.id, target))

		case wire.MsgError:
			code, msg, _ := wire.DecodeError(fr.Payload)
			log.Printf("zone %d error: code=%d msg=%q", zl.id, code, msg)
		}
	}
}

func (s *Server) tryCommitOnAttachAck(ackZoneID uint32) {
	var commits []struct{ sid shared.SessionID; from shared.ZoneID }
	s.xferMu.Lock()
	for sid, xs := range s.inflight {
		if uint32(xs.To) == ackZoneID {
			commits = append(commits, struct{sid shared.SessionID; from shared.ZoneID}{sid, xs.From})
			delete(s.inflight, sid)
		}
	}
	s.xferMu.Unlock()

	for _, c := range commits {
		_ = s.zoneSend(uint32(c.from), wire.MsgTransferCommit, wire.EncodeTransferCommit(c.sid))
		if st, ok := s.getBySID(c.sid); ok {
			s.sendReliableText(st, "XFER_COMMIT")
		}
	}
}

// ---- tiny little-endian helpers (avoid extra deps) ----
func binaryLEU16(b []byte) uint16 { return uint16(b[0]) | uint16(b[1])<<8 }
func binaryLEU32(b []byte) uint32 { return uint32(binaryLEU16(b[0:2])) | uint32(binaryLEU16(b[2:4]))<<16 }
func binaryLEU64(b []byte) uint64 { return uint64(binaryLEU32(b[0:4])) | uint64(binaryLEU32(b[4:8]))<<32 }
