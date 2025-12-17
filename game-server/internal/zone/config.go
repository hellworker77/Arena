package zone

import "game-server/internal/persist"

type Config struct {
	ListenAddr string
	ZoneID     uint32
	TickHz     int

	AOIRadius int16
	CellSize  int16

	MaxMoveEvents  int
	MaxStateEvents int
	MaxEventEvents int

	BudgetBytes int
	StateEveryTicks int

	// persistence
	SaveEveryTicks int
	Store persist.Store
	SaveQ  *persist.SaveQueue
}
