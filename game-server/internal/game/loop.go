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
    return &Loop{
        engine:   engine,
        tickRate: tickRate,
        cmds:     make(chan command, 4096),
    }
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
func (l *Loop) Run(ctx context.Context, onSnapshot func(snapshotPayload []byte)) {
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
            if onSnapshot != nil {
                onSnapshot(l.engine.BuildPlayersSnapshot())
            }
        }
    }
}
