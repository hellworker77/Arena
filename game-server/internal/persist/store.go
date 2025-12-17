package persist

import (
	"context"

	"game-server/internal/shared"
)

type CharacterState struct {
	CharacterID shared.CharacterID `json:"character_id"`
	ZoneID      shared.ZoneID      `json:"zone_id"`
	X           int16              `json:"x"`
	Y           int16              `json:"y"`
	HP          uint16             `json:"hp"`
	ServerTick  uint32             `json:"server_tick"`
}

type Store interface {
	LoadCharacter(ctx context.Context, id shared.CharacterID) (CharacterState, bool, error)
	SaveCharacter(ctx context.Context, st CharacterState) error
}
