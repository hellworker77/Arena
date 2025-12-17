package persist

import (
	"context"
	"sync"
	"time"

	"game-server/internal/shared"
)

type SaveQueue struct {
	store Store

	mu sync.Mutex
	pending map[shared.CharacterID]CharacterState
	order []shared.CharacterID
	maxPending int
	wake chan struct{}
}

func NewSaveQueue(store Store, maxPending int) *SaveQueue {
	if maxPending <= 0 { maxPending = 10_000 }
	return &SaveQueue{
		store: store,
		pending: make(map[shared.CharacterID]CharacterState),
		maxPending: maxPending,
		wake: make(chan struct{}, 1),
	}
}

func (q *SaveQueue) Enqueue(st CharacterState) {
	q.mu.Lock()
	_, exists := q.pending[st.CharacterID]
	q.pending[st.CharacterID] = st
	if !exists {
		q.order = append(q.order, st.CharacterID)
		for len(q.order) > q.maxPending {
			old := q.order[0]
			q.order = q.order[1:]
			delete(q.pending, old)
		}
	}
	q.mu.Unlock()
	select { case q.wake <- struct{}{}: default: }
}

func (q *SaveQueue) Run(ctx context.Context) error {
	t := time.NewTicker(250 * time.Millisecond)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			q.flushSome(ctx, 10_000)
			return nil
		case <-q.wake:
			q.flushSome(ctx, 256)
		case <-t.C:
			q.flushSome(ctx, 256)
		}
	}
}

func (q *SaveQueue) flushSome(ctx context.Context, n int) {
	if n <= 0 { return }
	var batch []CharacterState
	q.mu.Lock()
	for len(batch) < n && len(q.order) > 0 {
		cid := q.order[0]
		q.order = q.order[1:]
		st, ok := q.pending[cid]
		if ok {
			delete(q.pending, cid)
			batch = append(batch, st)
		}
	}
	q.mu.Unlock()

	for _, st := range batch {
		_ = q.store.SaveCharacter(ctx, st)
	}
}
