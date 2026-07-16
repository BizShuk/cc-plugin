package model

import (
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestStateStoreCursorPosition(t *testing.T) {
	originalStatePath := viper.GetString("state.db_path")
	t.Cleanup(func() {
		viper.Set("state.db_path", originalStatePath)
	})
	viper.Set("state.db_path", filepath.Join(t.TempDir(), "state.db"))

	store, err := NewStateStore()
	if err != nil {
		t.Fatalf("open state store: %v", err)
	}
	t.Cleanup(func() {
		if err := store.Close(); err != nil {
			t.Errorf("close state store: %v", err)
		}
	})

	want := CursorPosition{LastTS: 100, LastID: 7}
	if err := store.SetCursorPosition("claude-mem-export", want); err != nil {
		t.Fatalf("set cursor position: %v", err)
	}

	got, err := store.GetCursorPosition("claude-mem-export")
	if err != nil {
		t.Fatalf("get cursor position: %v", err)
	}
	if got != want {
		t.Fatalf("expected cursor position %+v, got %+v", want, got)
	}
}
