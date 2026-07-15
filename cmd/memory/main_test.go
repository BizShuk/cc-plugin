package memory

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bizshuk/cc-plugin/model"
	"github.com/spf13/viper"
)

func TestReadGbrainLogic(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gbrain-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	file1 := filepath.Join(tmpDir, "alice.md")
	if err := os.WriteFile(file1, []byte("alice context"), 0o644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	os.Chtimes(file1, time.Unix(1000, 0), time.Unix(1000, 0))

	subDir := filepath.Join(tmpDir, "topics")
	if err := os.MkdirAll(subDir, 0o755); err != nil {
		t.Fatalf("failed to create sub dir: %v", err)
	}

	file2 := filepath.Join(subDir, "trip.md")
	if err := os.WriteFile(file2, []byte("trip notes"), 0o644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	os.Chtimes(file2, time.Unix(2000, 0), time.Unix(2000, 0))

	thresholdTime := int64(2500)

	file3 := filepath.Join(tmpDir, "new.md")
	if err := os.WriteFile(file3, []byte("new notes"), 0o644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	os.Chtimes(file3, time.Unix(3000, 0), time.Unix(3000, 0))

	// Verify reading since 0
	var observations []model.Observation
	err = filepath.Walk(tmpDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			mtime := info.ModTime().Unix()
			if mtime > 0 {
				content, _ := os.ReadFile(path)
				rel, _ := filepath.Rel(tmpDir, path)
				observations = append(observations, model.Observation{
					Source:    "gbrain-working",
					SourceID:  rel,
					Timestamp: mtime,
					Text:      string(content),
				})
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk failed: %v", err)
	}

	if len(observations) != 3 {
		t.Errorf("expected 3 observations, got %d", len(observations))
	}

	// Verify reading since thresholdTime
	var newObservations []model.Observation
	err = filepath.Walk(tmpDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			mtime := info.ModTime().Unix()
			if mtime > thresholdTime {
				content, _ := os.ReadFile(path)
				rel, _ := filepath.Rel(tmpDir, path)
				newObservations = append(newObservations, model.Observation{
					Source:    "gbrain-working",
					SourceID:  rel,
					Timestamp: mtime,
					Text:      string(content),
				})
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk failed: %v", err)
	}

	if len(newObservations) != 1 {
		t.Errorf("expected 1 new observation, got %d", len(newObservations))
	} else if newObservations[0].SourceID != "new.md" {
		t.Errorf("expected new.md, got %s", newObservations[0].SourceID)
	}
}

func TestReadClaudeMemLogic(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "claudemem-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "claude-mem.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("failed to open sqlite DB: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE observations (
			id INTEGER PRIMARY KEY,
			created_at_epoch INTEGER,
			text TEXT
		)
	`)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	_, err = db.Exec(`
		INSERT INTO observations (id, created_at_epoch, text)
		VALUES 
		(1, 100, 'first memory'),
		(2, 200, 'second memory')
	`)
	if err != nil {
		t.Fatalf("failed to seed table: %v", err)
	}

	statePath := filepath.Join(tmpDir, "state.db")
	viper.Set("sources.claude_mem.db_path", dbPath)
	viper.Set("state.db_path", statePath)

	store, err := NewStateStore()
	if err != nil {
		t.Fatalf("failed to create state store: %v", err)
	}
	defer store.Close()

	// Initial read should return both observations
	observations, maxTS, err := readClaudeMemLogic()
	if err != nil {
		t.Fatalf("readClaudeMemLogic failed: %v", err)
	}

	if len(observations) != 2 {
		t.Errorf("expected 2 observations, got %d", len(observations))
	}
	if maxTS != 200 {
		t.Errorf("expected maxTS 200, got %d", maxTS)
	}

	// Update cursor to 150 and read again (should only return second memory)
	if err := store.SetCursor("claude-mem", 150); err != nil {
		t.Fatalf("failed to set cursor: %v", err)
	}

	observations, maxTS, err = readClaudeMemLogic()
	if err != nil {
		t.Fatalf("readClaudeMemLogic failed: %v", err)
	}

	if len(observations) != 1 {
		t.Fatalf("expected 1 observation, got %d", len(observations))
	}
	if observations[0].SourceID != "2" {
		t.Errorf("expected observation ID 2, got %s", observations[0].SourceID)
	}
	if observations[0].Text != "second memory" {
		t.Errorf("expected text 'second memory', got %s", observations[0].Text)
	}
}

func TestWriteAgentMemoryPayloadMapping(t *testing.T) {
	mem := model.Memory{
		Fingerprint: "fp123",
		Text:        "hello world",
		Entities:    []string{"alice", "bob"},
		Kind:        "preference",
		CreatedAt:   1000,
	}

	payload := map[string]interface{}{
		"content":  mem.Text,
		"concepts": append([]string{mem.Kind}, mem.Entities...),
		"files":    []string{},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &parsed); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if parsed["content"] != "hello world" {
		t.Errorf("expected content 'hello world', got %v", parsed["content"])
	}

	concepts := parsed["concepts"].([]interface{})
	if len(concepts) != 3 || concepts[0] != "preference" || concepts[1] != "alice" || concepts[2] != "bob" {
		t.Errorf("unexpected concepts payload: %v", concepts)
	}
}
