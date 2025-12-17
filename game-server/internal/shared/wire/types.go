package wire

// Internal Gateway <-> Zone link protocol.
// Frames on TCP: [len:uint32][type:uint8][payload...]

type MsgType uint8

const (
	// Gateway -> Zone
	MsgAttachPlayer MsgType = 1
	MsgDetachPlayer MsgType = 2
	MsgPlayerInput  MsgType = 3

	// Zone -> Gateway
	MsgAttachAck MsgType = 101
	MsgError     MsgType = 102
	MsgSnapshot  MsgType = 103
)

type ErrCode uint16

const (
	ErrUnknown  ErrCode = 0
	ErrBadMsg   ErrCode = 1
	ErrNoPlayer ErrCode = 2
)
