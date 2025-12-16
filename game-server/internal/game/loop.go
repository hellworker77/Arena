package game

import (
	"context"
	"game-server/pkg/protocol"
	"time"
)

// Loop owns the authoritative simulation thread.
// All Engine mutations happen only inside Run().
type Loop struct {
	engine   *Engine
	tickRate int
	cmds     chan command

	// replication settings
	interestRadius    float32
	fullEveryTicks    uint32
	gridCellSize      float32
	maxSnapshotBytes  int
	serverTickCounter uint32
}

// SnapshotFrame contains per-client snapshot payloads for a single server tick.
// The UDP layer can broadcast these payloads without touching the world state.
type SnapshotFrame struct {
	ServerTick uint32
	Payloads   map[string][]byte // clientKey -> payload
}

type command interface {
	apply(e *Engine)
}

type cmdFn func(*Engine)

func (f cmdFn) apply(e *Engine) { f(e) }

func NewLoop(engine *Engine, tickRate int) *Loop {
	if tickRate <= 0 {
		tickRate = 20
	}
	l := &Loop{
		engine:         engine,
		tickRate:       tickRate,
		cmds:           make(chan command, 4096),
		interestRadius: 25,
		fullEveryTicks: 20, // 1s at 20Hz
		gridCellSize:   8,
		maxSnapshotBytes: 1200,
	}
	// Apply defaults to engine.
	l.engine.ConfigureSpatialCellSize(l.gridCellSize)
	l.engine.ConfigureReplicationBudget(l.maxSnapshotBytes)
	return l
}

// ConfigureReplication overrides relevance/delta settings.
func (l *Loop) ConfigureReplication(interestRadius float32, fullEveryTicks uint32) {
	if interestRadius > 0 {
		l.interestRadius = interestRadius
		// If caller hasn't set a custom grid cell size, keep it at or below interest radius.
		if l.gridCellSize <= 0 || l.gridCellSize > interestRadius {
			l.gridCellSize = interestRadius
			l.engine.ConfigureSpatialCellSize(l.gridCellSize)
		}
	}
	if fullEveryTicks > 0 {
		l.fullEveryTicks = fullEveryTicks
	}
}

// ConfigureSpatialGrid sets the cell size for interest queries. Smaller = more cells, less contention per query.
func (l *Loop) ConfigureSpatialGrid(cellSize float32) {
	if cellSize <= 0 {
		return
	}
	l.gridCellSize = cellSize
	l.engine.ConfigureSpatialCellSize(cellSize)
}

// ConfigureSnapshotBudget sets a hard byte budget for snapshot payloads per client.
func (l *Loop) ConfigureSnapshotBudget(maxBytes int) {
	l.maxSnapshotBytes = maxBytes
	l.engine.ConfigureReplicationBudget(maxBytes)
}

// AddPlayer enqueues player creation and waits for the entity id.
func (l *Loop) AddPlayer(clientKey string) uint32 {
	resp := make(chan uint32, 1)
	l.cmds <- cmdFn(func(e *Engine) {
		id := e.AddPlayer(clientKey)
		resp <- uint32(id)
	})
	return <-resp
}

// RemovePlayer enqueues player removal.
func (l *Loop) RemovePlayer(clientKey string) {
	l.cmds <- cmdFn(func(e *Engine) { e.RemovePlayer(clientKey) })
}

// QueueInput enqueues input for processing on the next tick.
func (l *Loop) QueueInput(clientKey string, in protocol.Input) {
	// Input is best-effort: if the loop is overloaded, drop rather than block networking.
	select {
	case l.cmds <- cmdFn(func(e *Engine) { _ = e.QueueInput(clientKey, in) }):
	default:
		// drop
	}
}

// Run starts the simulation loop. It blocks until ctx is cancelled.
// onSnapshot is called once per tick with a fresh snapshot payload.
func (l *Loop) Run(ctx context.Context, onFrame func(frame SnapshotFrame)) {
	tick := time.NewTicker(time.Second / time.Duration(l.tickRate))
	defer tick.Stop()

	dt := float32(1.0 / float32(l.tickRate))

	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			// Drain commands quickly.
			for {
				select {
				case c := <-l.cmds:
					c.apply(l.engine)
				default:
					goto STEP
				}
			}
		STEP:
			l.engine.Step(dt)
			l.serverTickCounter++
			if onFrame != nil {
				frame := SnapshotFrame{ServerTick: l.serverTickCounter, Payloads: make(map[string][]byte, 16)}
				// Build per-client snapshots inside the simulation thread.
				for clientKey := range l.engine.players {
					if payload, ok := l.engine.BuildSnapshotForClient(clientKey, l.serverTickCounter, l.interestRadius, l.fullEveryTicks); ok {
						frame.Payloads[clientKey] = payload
					}
				}
				onFrame(frame)
			}
		}
	}
}
