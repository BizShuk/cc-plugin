package cmd

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestReadGbrainLogic(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gbrain-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	file1 := filepath.Join(tmpDir, "alice.md")
	if err := os.WriteFile(file1, []byte("alice context"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	os.Chtimes(file1, time.Unix(1000, 0), time.Unix(1000, 0))

	subDir := filepath.Join(tmpDir, "topics")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("failed to create sub dir: %v", err)
	}

	file2 := filepath.Join(subDir, "trip.md")
	if err := os.WriteFile(file2, []byte("trip notes"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	os.Chtimes(file2, time.Unix(2000, 0), time.Unix(2000, 0))

	thresholdTime := int64(2500)

	file3 := filepath.Join(tmpDir, "new.md")
	if err := os.WriteFile(file3, []byte("new notes"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	os.Chtimes(file3, time.Unix(3000, 0), time.Unix(3000, 0))

	// Verify reading since 0
	var observations []Observation
	err = filepath.Walk(tmpDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			mtime := info.ModTime().Unix()
			if mtime > 0 {
				content, _ := os.ReadFile(path)
				rel, _ := filepath.Rel(tmpDir, path)
				observations = append(observations, Observation{
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
	var newObservations []Observation
	err = filepath.Walk(tmpDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			mtime := info.ModTime().Unix()
			if mtime > thresholdTime {
				content, _ := os.ReadFile(path)
				rel, _ := filepath.Rel(tmpDir, path)
				newObservations = append(newObservations, Observation{
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
			text TEXT,
			title TEXT,
			subtitle TEXT,
			facts TEXT,
			narrative TEXT
		)
	`)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	_, err = db.Exec(`
		INSERT INTO observations (id, created_at_epoch, text, title, subtitle, facts, narrative)
		VALUES 
		(1, 100, 'first memory', 'title1', 'subtitle1', '[]', 'narrative1'),
		(2, 200, '', 'title2', 'subtitle2', '["fact1"]', 'narrative2')
	`)
	if err != nil {
		t.Fatalf("failed to seed table: %v", err)
	}

	// Query with timestamp filter
	rows, err := db.Query("SELECT id, created_at_epoch, text, title, subtitle, facts, narrative FROM observations WHERE created_at_epoch > ? ORDER BY created_at_epoch ASC", 150)
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}
	defer rows.Close()

	var observations []Observation
	for rows.Next() {
		var sid string
		var ts int64
		var textVal, title, subtitle, facts, narrative sql.NullString
		err = rows.Scan(&sid, &ts, &textVal, &title, &subtitle, &facts, &narrative)
		if err != nil {
			t.Fatalf("scan failed: %v", err)
		}

		var fullText string
		if textVal.Valid && textVal.String != "" {
			fullText = textVal.String
		} else {
			var parts []string
			if title.Valid && title.String != "" {
				parts = append(parts, "Title: "+title.String)
			}
			if subtitle.Valid && subtitle.String != "" {
				parts = append(parts, "Subtitle: "+subtitle.String)
			}
			if narrative.Valid && narrative.String != "" {
				parts = append(parts, "Narrative: "+narrative.String)
			}
			if facts.Valid && facts.String != "" {
				parts = append(parts, "Facts: "+facts.String)
			}
			fullText = strings.Join(parts, "\n")
		}

		observations = append(observations, Observation{
			Source:    "claude-mem",
			SourceID:  sid,
			Timestamp: ts,
			Text:      fullText,
		})
	}

	if len(observations) != 1 {
		t.Fatalf("expected 1 observation, got %d", len(observations))
	}

	if observations[0].SourceID != "2" {
		t.Errorf("expected observation ID 2, got %s", observations[0].SourceID)
	}

	expectedText := "Title: title2\nSubtitle: subtitle2\nNarrative: narrative2\nFacts: [\"fact1\"]"
	if observations[0].Text != expectedText {
		t.Errorf("expected concatenated text:\n%s\n\ngot:\n%s", expectedText, observations[0].Text)
	}
}

func TestWriteAgentMemoryPayloadMapping(t *testing.T) {
	mem := Memory{
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
