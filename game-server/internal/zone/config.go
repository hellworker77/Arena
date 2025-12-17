package zone

import "game-server/internal/persist"

type Config struct {
	ListenAddr string
	HTTPAddr   string

	ZoneID     uint32
	TickHz     int

	AOIRadius int16
	CellSize  int16

	BudgetBytes int
	StateEveryTicks int

	// persistence
	SaveEveryTicks int
	Store persist.Store
	SaveQ  *persist.SaveQueue

	// snapshots
	SnapshotEveryTicks int
	SnapshotStore persist.SnapshotStore
	SnapshotQ *persist.SnapshotQueue

	// AI budget
	AIBudgetPerTick int

	// transfer (boundary-based toy policy)
	TransferTargetZone uint32
	TransferBoundaryX int16
	TransferTimeoutTicks uint32
}
