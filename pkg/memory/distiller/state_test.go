package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStateStore(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "distiller-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "state.db")
	store, err := NewStateStore(dbPath)
	if err != nil {
		t.Fatalf("failed to create state store: %v", err)
	}
	defer store.Close()

	// Test Cursor
	val, err := store.GetCursor("test-source")
	if err != nil {
		t.Errorf("GetCursor failed: %v", err)
	}
	if val != 0 {
		t.Errorf("expected default cursor 0, got %d", val)
	}

	err = store.SetCursor("test-source", 12345)
	if err != nil {
		t.Errorf("SetCursor failed: %v", err)
	}

	val, err = store.GetCursor("test-source")
	if err != nil {
		t.Errorf("GetCursor failed: %v", err)
	}
	if val != 12345 {
		t.Errorf("expected cursor 12345, got %d", val)
	}

	// Test Seen
	seenCount, err := store.RecordSeen("fp1", "sourceA")
	if err != nil {
		t.Errorf("RecordSeen failed: %v", err)
	}
	if seenCount != 1 {
		t.Errorf("expected seen count 1, got %d", seenCount)
	}

	seenCount, err = store.RecordSeen("fp1", "sourceB")
	if err != nil {
		t.Errorf("RecordSeen failed: %v", err)
	}
	if seenCount != 2 {
		t.Errorf("expected seen count 2, got %d", seenCount)
	}

	// Test Distilled
	distilled, err := store.AlreadyDistilled("sourceA", "id1")
	if err != nil {
		t.Errorf("AlreadyDistilled failed: %v", err)
	}
	if distilled {
		t.Errorf("expected not distilled, but got true")
	}

	err = store.MarkDistilled("sourceA", "id1", 1000)
	if err != nil {
		t.Errorf("MarkDistilled failed: %v", err)
	}

	distilled, err = store.AlreadyDistilled("sourceA", "id1")
	if err != nil {
		t.Errorf("AlreadyDistilled failed: %v", err)
	}
	if !distilled {
		t.Errorf("expected distilled, but got false")
	}

	// Test Prune
	err = store.MarkDistilled("sourceA", "id2", 500)
	if err != nil {
		t.Errorf("MarkDistilled failed: %v", err)
	}

	items, err := store.DueForPrune(800)
	if err != nil {
		t.Errorf("DueForPrune failed: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("expected 1 item due for prune, got %d", len(items))
	} else {
		if items[0].SourceID != "id2" {
			t.Errorf("expected item id2, got %s", items[0].SourceID)
		}
	}

	err = store.DropDistilled("sourceA", "id2")
	if err != nil {
		t.Errorf("DropDistilled failed: %v", err)
	}

	items, err = store.DueForPrune(800)
	if err != nil {
		t.Errorf("DueForPrune failed: %v", err)
	}
	if len(items) != 0 {
		t.Errorf("expected 0 items due for prune after drop, got %d", len(items))
	}
}
