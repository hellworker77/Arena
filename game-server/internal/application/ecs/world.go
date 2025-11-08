package ecs

import (
	"game-server/internal/application/ecs/ecs_signatures/runtime"
	"game-server/internal/application/ecs/ecs_signatures/static"
)

type EntityRecord struct {
	Archetype *Archetype
	Index     int
}

type World struct {
	archetypes []*Archetype
	entities   map[static.EntityID]EntityRecord
	nextID     static.EntityID
	Grid       *HierarchicalGrid
}

func NewWorld() *World {
	return &World{
		entities: make(map[static.EntityID]EntityRecord),
		Grid:     NewHierarchicalGrid(2.0, 10.0, 50.0),
	}
}

func (w *World) getOrCreateArchetype(sig Signature) *Archetype {
	for _, a := range w.archetypes {
		if a.Signature == sig {
			return a
		}
	}

	a := &Archetype{Signature: sig}
	w.archetypes = append(w.archetypes, a)
	return a
}

func copyIf[T any](fromSig, mask Signature, from, to []T, oldIndex, newIndex int) {
	if fromSig&mask != 0 {
		to[newIndex] = from[oldIndex]
	}
}

func (w *World) MoveEntity(eID static.EntityID, newSig Signature) {
	rec, ok := w.entities[eID]
	if !ok {
		return
	}

	oldA := rec.Archetype
	oldIndex := rec.Index

	newA := w.getOrCreateArchetype(newSig)
	newIndex := newA.InsertEmpty(eID)

	sig := oldA.Signature

	copyIf(sig, CCollider, oldA.Colliders, newA.Colliders, oldIndex, newIndex)
	copyIf(sig, CCombatAttrs, oldA.CombatAttrs, newA.CombatAttrs, oldIndex, newIndex)
	copyIf(sig, CEnemyPreset, oldA.EnemyPresets, newA.EnemyPresets, oldIndex, newIndex)
	copyIf(sig, CMovementAttrs, oldA.MovementAttrs, newA.MovementAttrs, oldIndex, newIndex)
	copyIf(sig, CProjectilePreset, oldA.ProjectilePresets, newA.ProjectilePresets, oldIndex, newIndex)
	copyIf(sig, CInventory, oldA.Inventories, newA.Inventories, oldIndex, newIndex)
	copyIf(sig, CWorldItem, oldA.WorldItems, newA.WorldItems, oldIndex, newIndex)

	copyIf(sig, CAttackCooldown, oldA.AttackCooldowns, newA.AttackCooldowns, oldIndex, newIndex)
	copyIf(sig, CExperience, oldA.Experiences, newA.Experiences, oldIndex, newIndex)
	copyIf(sig, CHealth, oldA.Healths, newA.Healths, oldIndex, newIndex)
	copyIf(sig, CLifespan, oldA.Lifespans, newA.Lifespans, oldIndex, newIndex)
	copyIf(sig, CPos, oldA.Positions, newA.Positions, oldIndex, newIndex)
	copyIf(sig, CProjectileState, oldA.ProjectileStates, newA.ProjectileStates, oldIndex, newIndex)
	copyIf(sig, CTarget, oldA.Targets, newA.Targets, oldIndex, newIndex)
	copyIf(sig, CVel, oldA.Velocities, newA.Velocities, oldIndex, newIndex)

	copyIf(sig, CPlayerTag, oldA.PlayerTags, newA.PlayerTags, oldIndex, newIndex)
	copyIf(sig, CEnemyTag, oldA.EnemyTags, newA.EnemyTags, oldIndex, newIndex)
	copyIf(sig, CNpcTag, oldA.NpcTags, newA.NpcTags, oldIndex, newIndex)
	copyIf(sig, CProjectileTag, oldA.ProjectileTags, newA.ProjectileTags, oldIndex, newIndex)
	copyIf(sig, CItemTag, oldA.ItemTags, newA.ItemTags, oldIndex, newIndex)

	oldHasPos := (oldA.Signature & CPos) != 0
	var oldPos runtime.Position
	if oldHasPos {
		oldPos = oldA.Positions[oldIndex]
	}

	moved := oldA.Remove(oldIndex)
	if moved != 0 {
		recMoved := w.entities[moved]
		recMoved.Index = oldIndex
		recMoved.Archetype = oldA
		w.entities[moved] = recMoved
	}

	w.entities[eID] = EntityRecord{Archetype: newA, Index: newIndex}

	newHasPos := (newSig & CPos) != 0

	switch {
	case oldHasPos && !newHasPos:
		w.Grid.Remove(eID, oldPos)
	case !oldHasPos && newHasPos:
		newPos := newA.Positions[newIndex]
		w.Grid.Insert(eID, newPos)
	}
}

func (w *World) GetEntity(eID static.EntityID) (EntityRecord, bool) {
	rec, ok := w.entities[eID]
	return rec, ok
}

func (w *World) CreateEntity(initialSig Signature) static.EntityID {
	eID := w.nextID
	w.nextID++

	a := w.getOrCreateArchetype(initialSig)
	index := a.InsertEmpty(eID)
	w.entities[eID] = EntityRecord{Archetype: a, Index: index}

	if initialSig&CPos != 0 {
		pos := a.Positions[index]
		w.Grid.Insert(eID, pos)
	}

	return eID
}

func (w *World) RemoveEntity(eID static.EntityID) {
	rec, ok := w.entities[eID]
	if !ok {
		return
	}

	moved := rec.Archetype.Remove(rec.Index)
	if moved != 0 {
		recMoved := w.entities[moved]
		recMoved.Index = rec.Index
		recMoved.Archetype = rec.Archetype
		w.entities[moved] = recMoved
	}

	delete(w.entities, eID)
}

func (w *World) HasComponent(entity static.EntityID, component Signature) bool {
	rec, ok := w.entities[entity]
	if !ok {
		return false
	}

	return rec.Archetype.Signature&component != 0
}

func Set[T any](w *World, eID static.EntityID, c Component[T], v T) {
	rec := w.entities[eID]

	w.MoveEntity(eID, rec.Archetype.Signature|c.Mask)
	rec = w.entities[eID]
	(*c.Slice(rec.Archetype))[rec.Index] = v
}

func Get[T any](w *World, eID static.EntityID, c Component[T]) *T {
	rec, ok := w.entities[eID]
	if !ok || rec.Archetype.Signature&c.Mask == 0 {
		return nil
	}
	return &(*c.Slice(rec.Archetype))[rec.Index]
}

func Remove[T any](w *World, eID static.EntityID, c Component[T]) {
	rec := w.entities[eID]
	w.MoveEntity(eID, rec.Archetype.Signature&^c.Mask)
}

// TODO: 6a	Rollback netcode (GGPO-подобное)	Для PvP / экшенов
// TODO: 6b	State Compression					Чтобы пакет был маленький
// TODO: 6c	Delta Snapshots						Передаём только изменения
// TODO: 6d	Interest Management					Не отправляем сущности вне видимости
// TODO: Good to have: Pooling for projectiles
