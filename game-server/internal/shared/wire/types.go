package wire

type MsgType uint8

const (
	// Gateway -> Zone
	MsgAttachPlayer        MsgType = 1
	MsgAttachWithState     MsgType = 4
	MsgDetachPlayer        MsgType = 2
	MsgPlayerInput         MsgType = 3

	// Zone -> Gateway
	MsgAttachAck           MsgType = 101
	MsgError               MsgType = 102
	MsgReplicate           MsgType = 103
	MsgTransfer            MsgType = 104
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
	RepSpawn   RepOp = 1
	RepDespawn RepOp = 2
	RepMove    RepOp = 3

	RepStateHP RepOp = 10
	RepEventText RepOp = 20
)
