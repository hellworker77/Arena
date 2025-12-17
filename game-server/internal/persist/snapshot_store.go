package persist

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Snapshot struct {
	ZoneID     uint32           `json:"zone_id"`
	ServerTick uint32           `json:"server_tick"`
	Entities   []SnapshotEntity `json:"entities"`
}

type SnapshotEntity struct {
	EID   uint32 `json:"eid"`
	Kind  uint8  `json:"kind"`
	Owner uint64 `json:"owner"`
	X     int16  `json:"x"`
	Y     int16  `json:"y"`
	VX    int16  `json:"vx"`
	VY    int16  `json:"vy"`
	HP    uint16 `json:"hp"`
}

type SnapshotStore interface {
	LoadSnapshot(ctx context.Context, zoneID uint32) (Snapshot, bool, error)
	SaveSnapshot(ctx context.Context, zoneID uint32, snap Snapshot) error
}

type JSONSnapshotStore struct{ Dir string }

func NewJSONSnapshotStore(dir string) (*JSONSnapshotStore, error) {
	if dir == "" {
		return nil, errors.New("snapshot dir required")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return &JSONSnapshotStore{Dir: dir}, nil
}

func (s *JSONSnapshotStore) path(zoneID uint32) string {
	return filepath.Join(s.Dir, "zone_"+itoaU64(uint64(zoneID))+"_snapshot.json")
}

func (s *JSONSnapshotStore) LoadSnapshot(ctx context.Context, zoneID uint32) (Snapshot, bool, error) {
	_ = ctx
	b, err := os.ReadFile(s.path(zoneID))
	if err != nil {
		if os.IsNotExist(err) {
			return Snapshot{}, false, nil
		}
		return Snapshot{}, false, err
	}
	var snap Snapshot
	if err := json.Unmarshal(b, &snap); err != nil {
		return Snapshot{}, false, err
	}
	return snap, true, nil
}

func (s *JSONSnapshotStore) SaveSnapshot(ctx context.Context, zoneID uint32, snap Snapshot) error {
	_ = ctx
	b, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path(zoneID) + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path(zoneID))
}
