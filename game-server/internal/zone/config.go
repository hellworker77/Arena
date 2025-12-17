package zone

import "game-server/internal/persist"

type Config struct {
	ListenAddr string
	ZoneID     uint32
	TickHz     int

	AOIRadius int16
	CellSize  int16

	BudgetBytes int
	StateEveryTicks int

	SaveEveryTicks int
	Store persist.Store
	SaveQ  *persist.SaveQueue

	// Transfer policy (toy boundary-based):
	// If TransferBoundaryX > 0 => transfer when X > boundary.
	// If TransferBoundaryX < 0 => transfer when X < boundary.
	TransferTargetZone uint32
	TransferBoundaryX int16
}
