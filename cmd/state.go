package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/bizshuk/cc-plugin/model"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type StateStore struct {
	db   *gorm.DB
	path string
}

func NewStateStore() (*StateStore, error) {
	path := expandPath(viper.GetString("state.db_path"))

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create state directory: %w", err)
	}

	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	store := &StateStore{
		db:   db,
		path: path,
	}

	if err := store.initSchema(); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *StateStore) initSchema() error {
	err := s.db.AutoMigrate(&model.Cursor{}, &model.Seen{}, &model.Distilled{})
	if err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}
	return nil
}

func (s *StateStore) GetCursor(source string) (int64, error) {
	var c model.Cursor
	err := s.db.First(&c, "source = ?", source).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to get cursor: %w", err)
	}
	return c.LastTs, nil
}

func (s *StateStore) SetCursor(source string, ts int64) error {
	c := model.Cursor{Source: source, LastTs: ts}
	err := s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "source"}},
		DoUpdates: clause.AssignmentColumns([]string{"last_ts"}),
	}).Create(&c).Error
	if err != nil {
		return fmt.Errorf("failed to set cursor: %w", err)
	}
	return nil
}

func (s *StateStore) RecordSeen(fingerprint string, source string) (int, error) {
	seen := model.Seen{
		Fingerprint: fingerprint,
		Source:      source,
		FirstSeen:   time.Now().Unix(),
	}
	err := s.db.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(&seen).Error
	if err != nil {
		return 0, fmt.Errorf("failed to record seen fingerprint: %w", err)
	}

	var count int64
	err = s.db.Model(&model.Seen{}).Where("fingerprint = ?", fingerprint).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count seen fingerprint: %w", err)
	}
	return int(count), nil
}

func (s *StateStore) AlreadyDistilled(source string, sourceID string) (bool, error) {
	var count int64
	err := s.db.Model(&model.Distilled{}).Where("source = ? AND source_id = ?", source, sourceID).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check already distilled: %w", err)
	}
	return count > 0, nil
}

func (s *StateStore) MarkDistilled(source string, sourceID string, at int64) error {
	d := model.Distilled{
		Source:      source,
		SourceID:    sourceID,
		DistilledAt: at,
	}
	err := s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "source"}, {Name: "source_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"distilled_at"}),
	}).Create(&d).Error
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
	var dbItems []model.Distilled
	err := s.db.Where("distilled_at < ?", beforeTS).Order("source, source_id").Find(&dbItems).Error
	if err != nil {
		return nil, fmt.Errorf("failed to query distilled for prune: %w", err)
	}

	items := make([]DistilledItem, len(dbItems))
	for i, d := range dbItems {
		items[i] = DistilledItem{
			Source:   d.Source,
			SourceID: d.SourceID,
		}
	}
	return items, nil
}

func (s *StateStore) DropDistilled(source string, sourceID string) error {
	err := s.db.Where("source = ? AND source_id = ?", source, sourceID).Delete(&model.Distilled{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete distilled: %w", err)
	}
	return nil
}

func (s *StateStore) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}
	return sqlDB.Close()
}

func Fingerprint(text string, entities []string) string {
	normalizedText := strings.ToLower(strings.Join(strings.Fields(text), " "))
	sortedEntities := make([]string, len(entities))
	copy(sortedEntities, entities)
	sort.Strings(sortedEntities)

	h := sha256.New()
	h.Write([]byte(normalizedText))
	h.Write([]byte("|"))
	h.Write([]byte(strings.Join(sortedEntities, "|")))
	return hex.EncodeToString(h.Sum(nil))
}

func expandPath(p string) string {
	if strings.HasPrefix(p, "~") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, p[1:])
	}
	return p
}
