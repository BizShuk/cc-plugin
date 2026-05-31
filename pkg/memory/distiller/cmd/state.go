package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Observation struct {
	Source    string `json:"source"`
	SourceID  string `json:"source_id"`
	Timestamp int64  `json:"timestamp"`
	Text      string `json:"text"`
}

type Memory struct {
	Fingerprint string   `json:"fingerprint"`
	Text        string   `json:"text"`
	Entities    []string `json:"entities"`
	Kind        string   `json:"kind"`
	CreatedAt   int64    `json:"created_at"`
}

type Fact struct {
	Fingerprint string     `json:"fingerprint"`
	Text        string     `json:"text"`
	Entities    []string   `json:"entities"`
	Evidence    [][]string `json:"evidence"`
	CreatedAt   int64      `json:"created_at"`
}

type StateStore struct {
	db   *sql.DB
	path string
}

func NewStateStore(path string) (*StateStore, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create state directory: %w", err)
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	store := &StateStore{
		db:   db,
		path: path,
	}

	if err := store.initSchema(); err != nil {
		db.Close()
		return nil, err
	}

	return store, nil
}

func (s *StateStore) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS cursor (
		source TEXT PRIMARY KEY,
		last_ts INTEGER NOT NULL
	);
	CREATE TABLE IF NOT EXISTS seen (
		fingerprint TEXT NOT NULL,
		source TEXT NOT NULL,
		first_seen INTEGER NOT NULL,
		PRIMARY KEY (fingerprint, source)
	);
	CREATE TABLE IF NOT EXISTS distilled (
		source TEXT NOT NULL,
		source_id TEXT NOT NULL,
		distilled_at INTEGER NOT NULL,
		PRIMARY KEY (source, source_id)
	);
	`
	_, err := s.db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}
	return nil
}

func (s *StateStore) GetCursor(source string) (int64, error) {
	var lastTS int64
	err := s.db.QueryRow("SELECT last_ts FROM cursor WHERE source = ?", source).Scan(&lastTS)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get cursor: %w", err)
	}
	return lastTS, nil
}

func (s *StateStore) SetCursor(source string, ts int64) error {
	query := `
	INSERT INTO cursor (source, last_ts) VALUES (?, ?)
	ON CONFLICT(source) DO UPDATE SET last_ts = excluded.last_ts
	`
	_, err := s.db.Exec(query, source, ts)
	if err != nil {
		return fmt.Errorf("failed to set cursor: %w", err)
	}
	return nil
}

func (s *StateStore) RecordSeen(fingerprint string, source string) (int, error) {
	_, err := s.db.Exec(
		"INSERT OR IGNORE INTO seen (fingerprint, source, first_seen) VALUES (?, ?, ?)",
		fingerprint, source, time.Now().Unix(),
	)
	if err != nil {
		return 0, fmt.Errorf("failed to record seen fingerprint: %w", err)
	}

	var count int
	err = s.db.QueryRow("SELECT COUNT(*) FROM seen WHERE fingerprint = ?", fingerprint).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count seen fingerprint: %w", err)
	}
	return count, nil
}

func (s *StateStore) AlreadyDistilled(source string, sourceID string) (bool, error) {
	var val int
	err := s.db.QueryRow(
		"SELECT 1 FROM distilled WHERE source = ? AND source_id = ?",
		source, sourceID,
	).Scan(&val)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check already distilled: %w", err)
	}
	return true, nil
}

func (s *StateStore) MarkDistilled(source string, sourceID string, at int64) error {
	query := `
	INSERT OR REPLACE INTO distilled (source, source_id, distilled_at) VALUES (?, ?, ?)
	`
	_, err := s.db.Exec(query, source, sourceID, at)
	if err != nil {
		return fmt.Errorf("failed to mark distilled: %w", err)
	}
	return nil
}

type DistilledItem struct {
	Source   string
	SourceID string
}

func (s *StateStore) DueForPrune(beforeTS int64) ([]DistilledItem, error) {
	rows, err := s.db.Query(
		"SELECT source, source_id FROM distilled WHERE distilled_at < ? ORDER BY source, source_id",
		beforeTS,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query distilled for prune: %w", err)
	}
	defer rows.Close()

	var items []DistilledItem
	for rows.Next() {
		var item DistilledItem
		if err := rows.Scan(&item.Source, &item.SourceID); err != nil {
			return nil, fmt.Errorf("failed to scan distilled item: %w", err)
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *StateStore) DropDistilled(source string, sourceID string) error {
	_, err := s.db.Exec("DELETE FROM distilled WHERE source = ? AND source_id = ?", source, sourceID)
	if err != nil {
		return fmt.Errorf("failed to delete distilled: %w", err)
	}
	return nil
}

func (s *StateStore) Close() error {
	return s.db.Close()
}
