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

	zoneMu sync.Mutex
	zoneW *bufio.Writer
	zoneR *bufio.Reader
	zoneC net.Conn

	sessionsMu sync.Mutex
	byRemote map[string]*sessionState
	bySID map[shared.SessionID]string
}

type sessionState struct {
	SID shared.SessionID
	CharID shared.CharacterID
	ZoneID shared.ZoneID
	LastHeard time.Time
}

func New(cfg Config) (*Server, error) {
	if cfg.UDPListenAddr=="" || cfg.ZoneTCPAddr=="" { return nil, errors.New("missing addresses") }
	if cfg.IdleTimeout<=0 { cfg.IdleTimeout = 30*time.Second }
	return &Server{
		cfg: cfg,
		byRemote: make(map[string]*sessionState),
		bySID: make(map[shared.SessionID]string),
	}, nil
}

func (s *Server) Start(ctx context.Context) error {
	udpAddr, err := net.ResolveUDPAddr("udp", s.cfg.UDPListenAddr)
	if err != nil { return err }
	s.udpConn, err = net.ListenUDP("udp", udpAddr)
	if err != nil { return err }
	defer s.udpConn.Close()

	if err := s.connectZone(); err != nil { return err }
	defer s.closeZone()

	go s.zoneReadLoop(ctx)
	go s.cleanupLoop(ctx)

	log.Printf("gateway up: udp=%s zone=%s", s.cfg.UDPListenAddr, s.cfg.ZoneTCPAddr)

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

func (s *Server) connectZone() error {
	c, err := net.Dial("tcp", s.cfg.ZoneTCPAddr)
	if err != nil { return err }
	s.zoneMu.Lock()
	s.zoneC = c
	s.zoneW = bufio.NewWriterSize(c, 64*1024)
	s.zoneR = bufio.NewReaderSize(c, 64*1024)
	s.zoneMu.Unlock()
	return nil
}

func (s *Server) closeZone() {
	s.zoneMu.Lock()
	defer s.zoneMu.Unlock()
	if s.zoneC != nil {
		_ = s.zoneC.Close()
		s.zoneC=nil; s.zoneW=nil; s.zoneR=nil
	}
}

func (s *Server) zoneSend(typ wire.MsgType, payload []byte) error {
	s.zoneMu.Lock(); defer s.zoneMu.Unlock()
	if s.zoneW==nil { return errors.New("zone link down") }
	return wire.WriteFrame(s.zoneW, typ, payload)
}

func (s *Server) handleUDPPacket(raddr *net.UDPAddr, b []byte) {
	line := string(b)
	ra := raddr.String()

	if len(line)>=5 && line[:5]=="HELLO" {
		var cid uint64
		if _, err := sscanf(line, "HELLO %d", &cid); err!=nil || cid==0 {
			_, _ = s.udpConn.WriteToUDP([]byte("ERR bad hello
"), raddr)
			return
		}
		st := s.ensureSession(ra, shared.CharacterID(cid))
		if err := s.zoneSend(wire.MsgAttachPlayer, wire.EncodeAttachPlayer(st.SID, st.CharID, st.ZoneID)); err != nil {
			_, _ = s.udpConn.WriteToUDP([]byte("ERR zone down
"), raddr); return
		}
		_, _ = s.udpConn.WriteToUDP([]byte("OK "+st.SID.String()+"
"), raddr)
		return
	}

	if len(line)>=2 && line[:2]=="IN" {
		var tick uint32
		var mx,my int16
		if _, err := sscanf(line, "IN %d %d %d", &tick, &mx, &my); err != nil {
			_, _ = s.udpConn.WriteToUDP([]byte("ERR bad input
"), raddr); return
		}
		st := s.getByRemote(ra)
		if st==nil {
			_, _ = s.udpConn.WriteToUDP([]byte("ERR no session
"), raddr); return
		}
		st.LastHeard = time.Now()
		_ = s.zoneSend(wire.MsgPlayerInput, wire.EncodePlayerInput(st.SID, tick, mx, my))
		return
	}

	_, _ = s.udpConn.WriteToUDP([]byte("ERR unknown
"), raddr)
}

func (s *Server) ensureSession(remote string, cid shared.CharacterID) *sessionState {
	s.sessionsMu.Lock(); defer s.sessionsMu.Unlock()
	if st, ok := s.byRemote[remote]; ok {
		st.LastHeard=time.Now(); return st
	}
	st := &sessionState{SID: shared.NewSessionID(), CharID: cid, ZoneID: 1, LastHeard: time.Now()}
	s.byRemote[remote]=st
	s.bySID[st.SID]=remote
	return st
}

func (s *Server) getByRemote(remote string) *sessionState {
	s.sessionsMu.Lock(); defer s.sessionsMu.Unlock()
	return s.byRemote[remote]
}

func (s *Server) getRemoteBySID(sid shared.SessionID) (string,bool) {
	s.sessionsMu.Lock(); defer s.sessionsMu.Unlock()
	ra, ok := s.bySID[sid]
	return ra, ok
}

func (s *Server) cleanupLoop(ctx context.Context) {
	t := time.NewTicker(2*time.Second); defer t.Stop()
	for {
		select { case <-ctx.Done(): return
		case <-t.C:
			now := time.Now()
			var toDetach []shared.SessionID
			s.sessionsMu.Lock()
			for ra, st := range s.byRemote {
				if now.Sub(st.LastHeard) > s.cfg.IdleTimeout {
					toDetach = append(toDetach, st.SID)
					delete(s.bySID, st.SID)
					delete(s.byRemote, ra)
				}
			}
			s.sessionsMu.Unlock()
			for _, sid := range toDetach {
				_ = s.zoneSend(wire.MsgDetachPlayer, wire.EncodeDetachPlayer(sid))
			}
		}
	}
}

func (s *Server) zoneReadLoop(ctx context.Context) {
	for {
		select { case <-ctx.Done(): return; default: }
		s.zoneMu.Lock(); r := s.zoneR; s.zoneMu.Unlock()
		if r==nil { time.Sleep(100*time.Millisecond); continue }
		fr, err := wire.ReadFrame(r)
		if err != nil { log.Printf("zone link read error: %v", err); return }
		switch fr.Type {
		case wire.MsgAttachAck:
		case wire.MsgReplicate:
			sid, _, ch, events, err := wire.DecodeReplicate(fr.Payload)
			if err != nil { log.Printf("bad replicate frame: %v", err); continue }
			remote, ok := s.getRemoteBySID(sid); if !ok { continue }
			raddr, err := net.ResolveUDPAddr("udp", remote); if err != nil { continue }
			switch ch {
			case wire.ChanMove:
				for _, ev := range events {
					switch ev.Op {
					case wire.RepSpawn:
						_, _ = s.udpConn.WriteToUDP([]byte(sprintf("SPAWN %d %d %d
", uint32(ev.EID), ev.X, ev.Y)), raddr)
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
					if ev.Op==wire.RepStateHP {
						_, _ = s.udpConn.WriteToUDP([]byte(sprintf("STAT %d hp=%d
", uint32(ev.EID), ev.Val)), raddr)
					}
				}
			case wire.ChanEvent:
				for _, ev := range events {
					if ev.Op==wire.RepEventText {
						_, _ = s.udpConn.WriteToUDP([]byte(sprintf("EV %s
", ev.Text)), raddr)
					}
				}
			}
		case wire.MsgError:
			code, msg, _ := wire.DecodeError(fr.Payload)
			log.Printf("zone error: code=%d msg=%q", code, msg)
		default:
			log.Printf("zone unknown msg type: %d", fr.Type)
		}
	}
}
