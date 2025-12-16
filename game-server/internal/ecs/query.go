package ecs

import (
	runtime2 "game-server/internal/ecs/ecs_signatures/runtime"
	static2 "game-server/internal/ecs/ecs_signatures/static"
	tag2 "game-server/internal/ecs/ecs_signatures/tag"
)

type Query struct {
	world *World
	mask  Signature
}

type QueryIter struct {
	query *Query

	aIdx int // archetype index
	eIdx int // entity index within the archetype
	curr *Archetype

	spatialFilter []static2.EntityID
	sfIdx         int
}

func (w *World) Query(mask Signature) Query {
	return Query{world: w, mask: mask}
}

func (q Query) RangeIter(center runtime2.Position, radius float32) *QueryIter {
	ids := q.world.Grid.Query(0, center, radius)

	return &QueryIter{
		query:         &q,
		spatialFilter: ids,
		sfIdx:         -1,
	}
}

func (q Query) Iter() *QueryIter {
	return &QueryIter{query: &q, aIdx: -1, eIdx: -1}
}

func (it *QueryIter) Next() bool {
	if it.spatialFilter != nil {
		for {
			it.sfIdx++
			if it.sfIdx >= len(it.spatialFilter) {
				return false
			}

			e := it.spatialFilter[it.sfIdx]
			rec, ok := it.query.world.entities[e]
			if !ok {
				continue
			}

			if rec.Archetype.Signature&it.query.mask != it.query.mask {
				continue
			}

			it.curr = rec.Archetype
			it.eIdx = rec.Index
			return true
		}
	}

	for {
		if it.curr == nil {
			it.aIdx++
			if it.aIdx >= len(it.query.world.archetypes) {
				return false
			}

			a := it.query.world.archetypes[it.aIdx]

			if a.Signature&it.query.mask != it.query.mask {
				continue
			}

			it.curr = a
			it.eIdx = -1
		}

		it.eIdx++

		if it.eIdx >= it.curr.Count {
			it.curr = nil
			continue
		}

		return true
	}
}

func (it *QueryIter) EntityID() static2.EntityID {
	return it.curr.EntityIDs[it.eIdx]
}

func (it *QueryIter) Position() *runtime2.Position {
	return &it.curr.Positions[it.eIdx]
}

func (it *QueryIter) Velocity() *runtime2.Velocity {
	return &it.curr.Velocities[it.eIdx]
}

func (it *QueryIter) Health() *runtime2.Health {
	return &it.curr.Healths[it.eIdx]
}

func (it *QueryIter) Experience() *runtime2.Experience {
	return &it.curr.Experiences[it.eIdx]
}

func (it *QueryIter) AttackCooldown() *runtime2.AttackCooldown {
	return &it.curr.AttackCooldowns[it.eIdx]
}

func (it *QueryIter) Lifespan() *runtime2.Lifespan {
	return &it.curr.Lifespans[it.eIdx]
}

func (it *QueryIter) Target() *runtime2.Target {
	return &it.curr.Targets[it.eIdx]
}

func (it *QueryIter) ProjectileState() *runtime2.ProjectileState {
	return &it.curr.ProjectileStates[it.eIdx]
}

func (it *QueryIter) Inventory() *runtime2.Inventory {
	return &it.curr.Inventories[it.eIdx]
}

func (it *QueryIter) WorldItem() *runtime2.WorldItem {
	return &it.curr.WorldItems[it.eIdx]
}

func (it *QueryIter) InterestState() *runtime2.InterestState {
	return &it.curr.InterestStates[it.eIdx]
}

func (it *QueryIter) Collider() *static2.Collider {
	return &it.curr.Colliders[it.eIdx]
}

func (it *QueryIter) CombatAttributes() *static2.CombatAttributes {
	return &it.curr.CombatAttrs[it.eIdx]
}

func (it *QueryIter) MovementAttributes() *static2.MovementAttributes {
	return &it.curr.MovementAttrs[it.eIdx]
}

func (it *QueryIter) EnemyPreset() *static2.EnemyPreset {
	return &it.curr.EnemyPresets[it.eIdx]
}

func (it *QueryIter) ProjectilePreset() *static2.ProjectilePreset {
	return &it.curr.ProjectilePresets[it.eIdx]
}

func (it *QueryIter) PlayerTag() *tag2.PlayerTag {
	return &it.curr.PlayerTags[it.eIdx]
}

func (it *QueryIter) EnemyTag() *tag2.EnemyTag {
	return &it.curr.EnemyTags[it.eIdx]
}

func (it *QueryIter) NpcTag() *tag2.NpcTag {
	return &it.curr.NpcTags[it.eIdx]
}

func (it *QueryIter) ProjectileTag() *tag2.ProjectileTag {
	return &it.curr.ProjectileTags[it.eIdx]
}

func (it *QueryIter) ItemTag() *tag2.ItemTag {
	return &it.curr.ItemTags[it.eIdx]
}
