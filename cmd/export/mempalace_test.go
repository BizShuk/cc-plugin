package export

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestMempalaceExport(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "mempalace-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "chroma.sqlite3")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test sqlite: %v", err)
	}

	// Create required schema tables
	err = db.Exec(`CREATE TABLE embeddings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		segment_id TEXT,
		embedding_id TEXT,
		seq_id BLOB,
		created_at TIMESTAMP
	)`).Error
	if err != nil {
		t.Fatalf("failed to create embeddings table: %v", err)
	}

	err = db.Exec(`CREATE TABLE embedding_metadata (
		id INTEGER,
		key TEXT,
		string_value TEXT,
		int_value INTEGER,
		float_value REAL,
		bool_value INTEGER
	)`).Error
	if err != nil {
		t.Fatalf("failed to create embedding_metadata table: %v", err)
	}

	// Insert test drawer 1
	err = db.Exec(`INSERT INTO embeddings (id, embedding_id) VALUES (1, 'drawer_test_1')`).Error
	if err != nil {
		t.Fatalf("failed to insert embedding: %v", err)
	}
	metadata := []struct {
		ID  int
		K   string
		Val string
	}{
		{1, "wing", "test-wing"},
		{1, "room", "test-room"},
		{1, "chroma:document", "Test document content"},
		{1, "source_file", "test.md"},
		{1, "filed_at", "2026-06-01T00:00:00Z"},
		{1, "added_by", "tester"},
	}
	for _, m := range metadata {
		err = db.Exec(`INSERT INTO embedding_metadata (id, key, string_value) VALUES (?, ?, ?)`, m.ID, m.K, m.Val).Error
		if err != nil {
			t.Fatalf("failed to insert metadata: %v", err)
		}
	}

	// Test ReadPalaceData
	rows, err := ReadPalaceData(dbPath)
	if err != nil {
		t.Fatalf("ReadPalaceData failed: %v", err)
	}
	if len(rows) != 1 {
		t.Errorf("expected 1 row, got %d", len(rows))
	}
	r := rows[0]
	if r.EmbeddingID != "drawer_test_1" || r.Wing != "test-wing" || r.Room != "test-room" || r.Document != "Test document content" {
		t.Errorf("unexpected row data: %+v", r)
	}

	// Test ExportCategories
	outDir := filepath.Join(tmpDir, "out")
	err = ExportCategories(rows, outDir)
	if err != nil {
		t.Fatalf("ExportCategories failed: %v", err)
	}
	csvPath := filepath.Join(outDir, "mempalace_categories.csv")
	csvFile, err := os.Open(csvPath)
	if err != nil {
		t.Fatalf("failed to open categories CSV: %v", err)
	}
	defer csvFile.Close()
	reader := csv.NewReader(csvFile)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("failed to read CSV: %v", err)
	}
	if len(records) != 2 || records[1][0] != "test-wing" || records[1][1] != "test-room" || records[1][2] != "drawer_test_1" {
		t.Errorf("unexpected CSV records: %+v", records)
	}

	// Test ExportData
	dataDir := filepath.Join(tmpDir, "data_out")
	err = ExportData(rows, dataDir)
	if err != nil {
		t.Fatalf("ExportData failed: %v", err)
	}
	mdPath := filepath.Join(dataDir, "test-wing", "test-room.md")
	mdBytes, err := os.ReadFile(mdPath)
	if err != nil {
		t.Fatalf("failed to read MD file: %v", err)
	}
	mdContent := string(mdBytes)
	if !strings.Contains(mdContent, "## drawer_test_1") || !strings.Contains(mdContent, "Test document content") {
		t.Errorf("unexpected MD content: %s", mdContent)
	}

	indexPath := filepath.Join(dataDir, "index.md")
	indexBytes, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("failed to read index.md: %v", err)
	}
	indexContent := string(indexBytes)
	if !strings.Contains(indexContent, "Palace Export") || !strings.Contains(indexContent, "test-wing") {
		t.Errorf("unexpected index.md content: %s", indexContent)
	}
}
