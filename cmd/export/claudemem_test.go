package export

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bizshuk/cc-plugin/model"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) {
	return 0, errors.New("write failed")
}

type writerFunc func([]byte) (int, error)

func (fn writerFunc) Write(p []byte) (int, error) {
	return fn(p)
}

func configureClaudeMemTestPaths(t *testing.T) (string, string) {
	t.Helper()

	tmpDir := t.TempDir()
	sourcePath := filepath.Join(tmpDir, "claude-mem.db")
	statePath := filepath.Join(tmpDir, "state.db")
	originalSourcePath := viper.GetString("sources.claude_mem.db_path")
	originalStatePath := viper.GetString("state.db_path")
	t.Cleanup(func() {
		viper.Set("sources.claude_mem.db_path", originalSourcePath)
		viper.Set("state.db_path", originalStatePath)
	})
	viper.Set("sources.claude_mem.db_path", sourcePath)
	viper.Set("state.db_path", statePath)

	return sourcePath, statePath
}

func createClaudeMemSource(t *testing.T, path string) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		t.Fatalf("open source database: %v", err)
	}
	if err := db.Exec(`
		CREATE TABLE observations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at_epoch INTEGER NOT NULL,
			text TEXT
		)
	`).Error; err != nil {
		t.Fatalf("create source schema: %v", err)
	}
	t.Cleanup(func() {
		sqlDB, err := db.DB()
		if err != nil {
			t.Errorf("get source sql database: %v", err)
			return
		}
		if err := sqlDB.Close(); err != nil {
			t.Errorf("close source database: %v", err)
		}
	})

	return db
}

func seedClaudeMemObservation(t *testing.T, db *gorm.DB, id, timestamp int64, text string) {
	t.Helper()
	if err := db.Exec(
		"INSERT INTO observations (id, created_at_epoch, text) VALUES (?, ?, ?)",
		id,
		timestamp,
		text,
	).Error; err != nil {
		t.Fatalf("seed source database: %v", err)
	}
}

func executeClaudeMemCommand(t *testing.T, args ...string) ([]model.Observation, error) {
	t.Helper()

	var output bytes.Buffer
	cmd := ClaudeMemCmd()
	cmd.SetOut(&output)
	cmd.SetArgs(args)
	cmd.SilenceErrors = true
	err := cmd.Execute()
	if err != nil {
		return nil, err
	}

	var observations []model.Observation
	if err := json.Unmarshal(output.Bytes(), &observations); err != nil {
		t.Fatalf("decode command output %q: %v", output.String(), err)
	}
	return observations, nil
}

func TestClaudeMemCmdDoesNotAdvanceCursorWhenOutputFails(t *testing.T) {
	sourcePath, _ := configureClaudeMemTestPaths(t)
	sourceDB := createClaudeMemSource(t, sourcePath)
	seedClaudeMemObservation(t, sourceDB, 1, 100, "memory")

	cmd := ClaudeMemCmd()
	cmd.SetOut(failingWriter{})
	cmd.SilenceErrors = true

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected output error, got nil")
	}
	if !strings.Contains(err.Error(), "write claude-mem export") {
		t.Fatalf("expected wrapped output error, got %v", err)
	}

	store, err := model.NewStateStore()
	if err != nil {
		t.Fatalf("open state store: %v", err)
	}
	t.Cleanup(func() {
		if err := store.Close(); err != nil {
			t.Errorf("close state store: %v", err)
		}
	})

	cursor, err := store.GetCursorPosition(claudeMemExportCursorSource)
	if err != nil {
		t.Fatalf("get claude-mem cursor: %v", err)
	}
	if cursor != (model.CursorPosition{}) {
		t.Fatalf("expected cursor to remain empty, got %+v", cursor)
	}
}

func TestClaudeMemCmdExportsLaterIDWithSameTimestamp(t *testing.T) {
	sourcePath, _ := configureClaudeMemTestPaths(t)
	sourceDB := createClaudeMemSource(t, sourcePath)
	seedClaudeMemObservation(t, sourceDB, 1, 100, "first")

	first, err := executeClaudeMemCommand(t)
	if err != nil {
		t.Fatalf("execute first export: %v", err)
	}
	if len(first) != 1 || first[0].SourceID != "1" {
		t.Fatalf("expected first export to contain ID 1, got %+v", first)
	}

	seedClaudeMemObservation(t, sourceDB, 2, 100, "second")
	second, err := executeClaudeMemCommand(t)
	if err != nil {
		t.Fatalf("execute second export: %v", err)
	}
	if len(second) != 1 || second[0].SourceID != "2" {
		t.Fatalf("expected second export to contain only ID 2, got %+v", second)
	}
}

func TestClaudeMemCmdDoesNotCreateMissingSourceDatabase(t *testing.T) {
	sourcePath, _ := configureClaudeMemTestPaths(t)

	_, err := executeClaudeMemCommand(t)
	if err == nil {
		t.Fatal("expected missing source database error, got nil")
	}
	if _, statErr := os.Stat(sourcePath); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("expected source database to remain absent, stat error: %v", statErr)
	}
}

func TestClaudeMemCmdReturnsCursorReadError(t *testing.T) {
	sourcePath, statePath := configureClaudeMemTestPaths(t)
	sourceDB := createClaudeMemSource(t, sourcePath)
	seedClaudeMemObservation(t, sourceDB, 1, 100, "memory")

	cmd := ClaudeMemCmd()
	cmd.SetArgs([]string{"--all"})
	cmd.SilenceErrors = true
	droppedCursorTable := false
	cmd.SetOut(writerFunc(func(p []byte) (int, error) {
		if droppedCursorTable {
			return len(p), nil
		}
		stateDB, err := gorm.Open(sqlite.Open(statePath), &gorm.Config{})
		if err != nil {
			return 0, err
		}
		if err := stateDB.Exec("DROP TABLE cursor").Error; err != nil {
			return 0, err
		}
		sqlDB, err := stateDB.DB()
		if err != nil {
			return 0, err
		}
		if err := sqlDB.Close(); err != nil {
			return 0, err
		}
		droppedCursorTable = true
		return len(p), nil
	}))

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected cursor read error, got nil")
	}
	if !strings.Contains(err.Error(), "get claude-mem export cursor") {
		t.Fatalf("expected cursor read error, got %v", err)
	}
}

func TestClaudeMemCmdWritesEmptyJSONAsArray(t *testing.T) {
	sourcePath, _ := configureClaudeMemTestPaths(t)
	createClaudeMemSource(t, sourcePath)

	var output bytes.Buffer
	cmd := ClaudeMemCmd()
	cmd.SetOut(&output)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute export: %v", err)
	}
	if output.String() != "[]\n" {
		t.Fatalf("expected empty JSON array, got %q", output.String())
	}
}
