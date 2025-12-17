package wire

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
)

var ErrFrameTooLarge = errors.New("wire: frame too large")

const MaxFrameSize = 1 << 20 // 1 MiB

type Frame struct {
	Type    MsgType
	Payload []byte
}

func WriteFrame(w *bufio.Writer, typ MsgType, payload []byte) error {
	l := 1 + len(payload)
	if l <= 0 || l > MaxFrameSize {
		return ErrFrameTooLarge
	}
	var hdr [4]byte
	binary.LittleEndian.PutUint32(hdr[:], uint32(l))
	if _, err := w.Write(hdr[:]); err != nil {
		return err
	}
	if err := w.WriteByte(byte(typ)); err != nil {
		return err
	}
	if len(payload) > 0 {
		if _, err := w.Write(payload); err != nil {
			return err
		}
	}
	return w.Flush()
}

func ReadFrame(r *bufio.Reader) (Frame, error) {
	var hdr [4]byte
	if _, err := io.ReadFull(r, hdr[:]); err != nil {
		return Frame{}, err
	}
	l := int(binary.LittleEndian.Uint32(hdr[:]))
	if l <= 0 || l > MaxFrameSize {
		return Frame{}, ErrFrameTooLarge
	}
	typB, err := r.ReadByte()
	if err != nil {
		return Frame{}, err
	}
	payloadLen := l - 1
	var payload []byte
	if payloadLen > 0 {
		payload = make([]byte, payloadLen)
		if _, err := io.ReadFull(r, payload); err != nil {
			return Frame{}, err
		}
	}
	return Frame{Type: MsgType(typB), Payload: payload}, nil
}
