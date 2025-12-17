package persist

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"game-server/internal/shared"
)

type JSONStore struct {
	Dir string
}

func NewJSONStore(dir string) (*JSONStore, error) {
	if dir == "" {
		return nil, errors.New("store dir required")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return &JSONStore{Dir: dir}, nil
}

func (s *JSONStore) path(id shared.CharacterID) string {
	return filepath.Join(s.Dir, "char_"+itoaU64(uint64(id))+".json")
}

func (s *JSONStore) LoadCharacter(ctx context.Context, id shared.CharacterID) (CharacterState, bool, error) {
	_ = ctx
	p := s.path(id)
	b, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return CharacterState{}, false, nil
		}
		return CharacterState{}, false, err
	}
	var st CharacterState
	if err := json.Unmarshal(b, &st); err != nil {
		return CharacterState{}, false, err
	}
	return st, true, nil
}

func (s *JSONStore) SaveCharacter(ctx context.Context, st CharacterState) error {
	_ = ctx
	b, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path(st.CharacterID) + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path(st.CharacterID))
}

// tiny local u64 itoa to avoid pulling strconv all over (but strconv is stdlib anyway)
func itoaU64(v uint64) string {
	if v == 0 {
		return "0"
	}
	var buf [32]byte
	i := len(buf)
	for v > 0 {
		i--
		buf[i] = byte('0' + (v % 10))
		v /= 10
	}
	return string(buf[i:])
}
