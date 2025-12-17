package zone

type Config struct {
	ListenAddr string
	ZoneID     uint32
	TickHz     int

	AOIRadius int16
	CellSize  int16

	// per-channel caps
	MaxMoveEvents  int
	MaxStateEvents int
	MaxEventEvents int

	// per-session total budget per tick (internal link payload bytes)
	BudgetBytes int

	// state sent every N ticks
	StateEveryTicks int
}
