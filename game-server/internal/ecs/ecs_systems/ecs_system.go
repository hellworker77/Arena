package ecs_systems

import (
	ecs2 "game-server/internal/ecs"
	"sync"
)

type ECSSystem interface {
	Run(w *ecs2.World, dt float32)
	Reads() ecs2.Signature
	Writes() ecs2.Signature
}

func RunSystems(w *ecs2.World, systems []ECSSystem, dt float32) {
	var wg sync.WaitGroup
	usedWrites := ecs2.Signature(0)

	for _, s := range systems {
		if usedWrites&s.Writes() != 0 {
			wg.Wait()
			usedWrites = 0
		}

		usedWrites |= s.Writes()

		wg.Add(1)
		go func(sys ECSSystem) {
			defer wg.Done()
			sys.Run(w, dt)
		}(s)
	}

	wg.Wait()
}
