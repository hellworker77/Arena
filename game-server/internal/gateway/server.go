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
}

type xferState struct {
	From shared.ZoneID
	To shared.ZoneID
	CID shared.CharacterID
	Interest wire.InterestMask
	X, Y int16
	HP uint16
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

	log.Printf("gateway up: udp=%s zones=%d proto=%d", s.cfg.UDPListenAddr, len(s.cfg.Zones), s.cfg.ProtoVersion)

	buf := make([]byte, 2048)
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

// UDP protocol (Step20 strict):
// - "HELLO <proto> <charID> [interestMask]"  (interest optional; default all)
// - "IN <tick> <mx> <my>"
// - "ACT <tick> <skill> <targetEID>"
func (s *Server) handleUDPPacket(raddr *net.UDPAddr, b []byte) {
	line := string(b)
	ra := raddr.String()

	if len(line) >= 5 && line[:5] == "HELLO" {
		var proto uint16
		var cid uint64
		var imask uint32
		n, err := sscanf(line, "HELLO %d %d %d", &proto, &cid, &imask)
		if err != nil || cid == 0 || proto == 0 {
			_, _ = s.udpConn.WriteToUDP([]byte("ERR bad hello
"), raddr)
			return
		}
		if proto != s.cfg.ProtoVersion {
			_, _ = s.udpConn.WriteToUDP([]byte("ERR bad proto
"), raddr)
			return
		}
		interest := wire.InterestMask(0)
		if n >= 3 {
			interest = wire.InterestMask(imask)
		}
		if interest == 0 {
			interest = wire.InterestMove | wire.InterestState | wire.InterestEvent | wire.InterestCombat
		}
		st := s.ensureSession(ra, proto, shared.CharacterID(cid), interest)
		if err := s.zoneSend(uint32(st.ZoneID), wire.MsgAttachPlayer, wire.EncodeAttachPlayer(st.SID, st.CharID, st.ZoneID, st.Interest)); err != nil {
			_, _ = s.udpConn.WriteToUDP([]byte("ERR zone down
"), raddr)
			return
		}
		_, _ = s.udpConn.WriteToUDP([]byte("OK "+st.SID.String()+"
"), raddr)
		return
	}

	if len(line) >= 2 && line[:2] == "IN" {
		var tick uint32
		var mx, my int16
		if _, err := sscanf(line, "IN %d %d %d", &tick, &mx, &my); err != nil {
			_, _ = s.udpConn.WriteToUDP([]byte("ERR bad input
"), raddr)
			return
		}
		st := s.getByRemote(ra)
		if st == nil {
			_, _ = s.udpConn.WriteToUDP([]byte("ERR no session
"), raddr)
			return
		}
		st.LastHeard = time.Now()
		_ = s.zoneSend(uint32(st.ZoneID), wire.MsgPlayerInput, wire.EncodePlayerInput(st.SID, tick, mx, my))
		return
	}

	if len(line) >= 3 && line[:3] == "ACT" {
		var tick uint32
		var skill uint16
		var target uint32
		if _, err := sscanf(line, "ACT %d %d %d", &tick, &skill, &target); err != nil {
			_, _ = s.udpConn.WriteToUDP([]byte("ERR bad act
"), raddr)
			return
		}
		st := s.getByRemote(ra)
		if st == nil {
			_, _ = s.udpConn.WriteToUDP([]byte("ERR no session
"), raddr)
			return
		}
		st.LastHeard = time.Now()
		_ = s.zoneSend(uint32(st.ZoneID), wire.MsgPlayerAction, wire.EncodePlayerAction(st.SID, tick, skill, shared.EntityID(target)))
		return
	}

	_, _ = s.udpConn.WriteToUDP([]byte("ERR unknown
"), raddr)
}

func (s *Server) ensureSession(remote string, proto uint16, cid shared.CharacterID, interest wire.InterestMask) *sessionState {
	s.sessionsMu.Lock()
	defer s.sessionsMu.Unlock()
	if st, ok := s.byRemote[remote]; ok {
		st.LastHeard = time.Now()
		return st
	}
	var min uint32 = 0
	for zid := range s.cfg.Zones {
		if min == 0 || zid < min { min = zid }
	}
	st := &sessionState{
		SID: shared.NewSessionID(),
		CharID: cid,
		ZoneID: shared.ZoneID(min),
		Interest: interest,
		LastHeard: time.Now(),
		Proto: proto,
	}
	s.byRemote[remote] = st
	s.bySID[st.SID] = remote
	return st
}

func (s *Server) getByRemote(remote string) *sessionState {
	s.sessionsMu.Lock(); defer s.sessionsMu.Unlock()
	return s.byRemote[remote]
}

func (s *Server) getRemoteBySID(sid shared.SessionID) (string, *sessionState, bool) {
	s.sessionsMu.Lock(); defer s.sessionsMu.Unlock()
	ra, ok := s.bySID[sid]
	if !ok { return "", nil, false }
	return ra, s.byRemote[ra], true
}

func (s *Server) setSessionZone(sid shared.SessionID, newZone shared.ZoneID) {
	s.sessionsMu.Lock(); defer s.sessionsMu.Unlock()
	ra, ok := s.bySID[sid]
	if !ok { return }
	if st := s.byRemote[ra]; st != nil { st.ZoneID = newZone }
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
				_ = s.zoneSend(uint32(d.zid), wire.MsgDetachPlayer, wire.EncodeDetachPlayer(d.sid))
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
			var abort []struct{
				sid shared.SessionID
				from shared.ZoneID
			}
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
			}
		}
	}
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
			// If this ACK corresponds to a transfer target, commit.
			// We key by SID: if inflight.To == this zone, commit old zone.
			// (Strict but simple: assumes only one inflight per SID.)
			// In real MMO you'd include transfer token.
			s.tryCommitOnAttachAck(zl.id)

		case wire.MsgReplicate:
			sid, _, ch, events, err := wire.DecodeReplicate(fr.Payload)
			if err != nil { continue }
			remote, _, ok := s.getRemoteBySID(sid)
			if !ok { continue }
			raddr, err := net.ResolveUDPAddr("udp", remote)
			if err != nil { continue }

			switch ch {
			case wire.ChanMove:
				for _, ev := range events {
					switch ev.Op {
					case wire.RepSpawn:
						_, _ = s.udpConn.WriteToUDP([]byte(sprintf("SPAWN %d %d %d kind=%d mask=%d
", uint32(ev.EID), ev.X, ev.Y, ev.Kind, uint32(ev.Mask))), raddr)
					case wire.RepDespawn:
						_, _ = s.udpConn.WriteToUDP([]byte(sprintf("DESPAWN %d
", uint32(ev.EID))), raddr)
					case wire.RepMove:
						_, _ = s.udpConn.WriteToUDP([]byte(sprintf("MOV %d %d %d
", uint32(ev.EID), ev.X, ev.Y)), raddr)
					}
				}
			case wire.ChanState:
				for _, ev := range events {
					if ev.Op == wire.RepStateHP {
						_, _ = s.udpConn.WriteToUDP([]byte(sprintf("STAT %d hp=%d
", uint32(ev.EID), ev.Val)), raddr)
					}
				}
			case wire.ChanEvent:
				for _, ev := range events {
					if ev.Op == wire.RepEventText {
						_, _ = s.udpConn.WriteToUDP([]byte(sprintf("EV %s
", ev.Text)), raddr)
					}
				}
			}

		case wire.MsgTransferPrepare:
			// Step13: 2PC transfer orchestration
			sid, cid, target, interest, x, y, hp, err := wire.DecodeTransferPrepare(fr.Payload)
			if err != nil { continue }

			_, st, ok := s.getRemoteBySID(sid)
			if !ok || st == nil { continue }

			// validate target exists
			s.zonesMu.Lock()
			_, okTarget := s.zones[uint32(target)]
			s.zonesMu.Unlock()
			if !okTarget {
				_ = s.zoneSend(uint32(st.ZoneID), wire.MsgTransferAbort, wire.EncodeTransferAbort(sid))
				continue
			}

			// register inflight
			s.xferMu.Lock()
			s.inflight[sid] = &xferState{
				From: st.ZoneID,
				To: shared.ZoneID(target),
				CID: cid,
				Interest: interest,
				X: x, Y: y, HP: hp,
				Started: time.Now(),
			}
			s.xferMu.Unlock()

			// attach with state to target
			_ = s.zoneSend(uint32(target), wire.MsgAttachWithState, wire.EncodeAttachWithState(sid, cid, shared.ZoneID(target), interest, x, y, hp))
			// route input immediately to target (client stays on gateway)
			s.setSessionZone(sid, shared.ZoneID(target))

			log.Printf("xfer prepare: sid=%s %d->%d", sid.String(), st.ZoneID, target)

		case wire.MsgError:
			code, msg, _ := wire.DecodeError(fr.Payload)
			log.Printf("zone %d error: code=%d msg=%q", zl.id, code, msg)
		}
	}
}

func (s *Server) tryCommitOnAttachAck(ackZoneID uint32) {
	// Find inflight transfers whose target == ackZoneID and commit them.
	var commits []struct{
		sid shared.SessionID
		from shared.ZoneID
	}
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
		log.Printf("xfer commit: sid=%s from=%d", c.sid.String(), c.from)
	}
}
