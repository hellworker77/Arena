package protocol

import (
    "bytes"
    "encoding/binary"
    "io"
)

// Input is a minimal gameplay input payload.
//
// It is intentionally small and versionable.
// Clients should send it in PTInput packets.
type Input struct {
    ClientTick uint32
    MoveX      float32 // -1..1
    MoveY      float32 // -1..1
    Speed      float32 // server can clamp/override
}

func MarshalInput(in Input) []byte {
    buf := new(bytes.Buffer)
    _ = binary.Write(buf, binary.LittleEndian, in.ClientTick)
    _ = binary.Write(buf, binary.LittleEndian, in.MoveX)
    _ = binary.Write(buf, binary.LittleEndian, in.MoveY)
    _ = binary.Write(buf, binary.LittleEndian, in.Speed)
    return buf.Bytes()
}

func UnmarshalInput(b []byte) (Input, error) {
    var in Input
    r := bytes.NewReader(b)
    if err := binary.Read(r, binary.LittleEndian, &in.ClientTick); err != nil {
        return in, err
    }
    if err := binary.Read(r, binary.LittleEndian, &in.MoveX); err != nil {
        return in, err
    }
    if err := binary.Read(r, binary.LittleEndian, &in.MoveY); err != nil {
        return in, err
    }
    if err := binary.Read(r, binary.LittleEndian, &in.Speed); err != nil {
        // Backward-compat: older clients may omit Speed.
        if err == io.EOF {
            in.Speed = 1
            return in, nil
        }
        return in, err
    }
    return in, nil
}
