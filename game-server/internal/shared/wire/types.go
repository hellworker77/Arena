package wire

// WireVersion is the contract for Gateway <-> Zone.
// Bump only with coordinated rollout.
const WireVersion uint16 = 1

type MsgType uint8

const (
	// Gateway -> Zone
	MsgAttachPlayer        MsgType = 1
	MsgAttachWithState     MsgType = 4
	MsgDetachPlayer        MsgType = 2
	MsgPlayerInput         MsgType = 3
	MsgPlayerAction        MsgType = 5

	// Transfer 2PC (Gateway -> Zone)
	MsgTransferCommit      MsgType = 6
	MsgTransferAbort       MsgType = 7

	// Zone -> Gateway
	MsgAttachAck           MsgType = 101
	MsgError               MsgType = 102
	MsgReplicate           MsgType = 103

	// Transfer 2PC (Zone -> Gateway)
	MsgTransferPrepare     MsgType = 104
)

type ErrCode uint16

const (
	ErrUnknown     ErrCode = 0
	ErrBadMsg      ErrCode = 1
	ErrNoPlayer    ErrCode = 2
	ErrBadAction   ErrCode = 3
	ErrCooldown    ErrCode = 4
	ErrOutOfRange  ErrCode = 5
	ErrTransfer    ErrCode = 6
)

type RepChannel uint8

const (
	ChanMove  RepChannel = 1
	ChanState RepChannel = 2
	ChanEvent RepChannel = 3
)

type RepOp uint8

const (
	RepSpawn      RepOp = 1
	RepDespawn    RepOp = 2
	RepMove       RepOp = 3

	RepStateHP    RepOp = 10
	RepEventText  RepOp = 20
)

// Interest layers (Step14)
type InterestMask uint32

const (
	InterestMove   InterestMask = 1 << 0
	InterestState  InterestMask = 1 << 1
	InterestEvent  InterestMask = 1 << 2
	InterestCombat InterestMask = 1 << 3
)

type EntityKind uint8

const (
	KindPlayer EntityKind = 1
	KindNPC    EntityKind = 2
)
