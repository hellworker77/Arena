package wire

// Internal Gateway <-> Zone link protocol.
// Frames on TCP: [len:uint32][type:uint8][payload...]
// Strict, no legacy.

type MsgType uint8

const (
	// Gateway -> Zone
	MsgAttachPlayer MsgType = 1
	MsgDetachPlayer MsgType = 2
	MsgPlayerInput  MsgType = 3

	// Zone -> Gateway
	MsgAttachAck MsgType = 101
	MsgError     MsgType = 102
	MsgReplicate MsgType = 103
)

type ErrCode uint16

const (
	ErrUnknown  ErrCode = 0
	ErrBadMsg   ErrCode = 1
	ErrNoPlayer ErrCode = 2
)

type RepChannel uint8

const (
	ChanMove  RepChannel = 1
	ChanState RepChannel = 2
	ChanEvent RepChannel = 3
)

type RepOp uint8

const (
	// movement channel ops
	RepSpawn   RepOp = 1
	RepDespawn RepOp = 2
	RepMove    RepOp = 3

	// state channel ops (toy)
	RepStateHP RepOp = 10

	// event channel ops (toy)
	RepEventText RepOp = 20
)
