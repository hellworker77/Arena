package ecs_systems

import (
	"game-server/internal/application/ecs"
	"sync"
)

type ECSSystem interface {
	Run(w *ecs.World, dt float32)
	Reads() ecs.Signature
	Writes() ecs.Signature
}

func RunSystems(w *ecs.World, systems []ECSSystem, dt float32) {
	var wg sync.WaitGroup
	usedWrites := ecs.Signature(0)

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
