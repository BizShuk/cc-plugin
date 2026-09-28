package model

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	gosdkconfig "github.com/bizshuk/gosdk/config"
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
	dbPath := viper.GetString("state.db_path")
	path := gosdkconfig.ExpandHome(dbPath)
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
	s := &StateStore{db: db, path: path}
	if err := s.initSchema(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *StateStore) initSchema() error {
	return s.db.AutoMigrate(&Cursor{})
}

// GetCursorPosition returns the timestamp and source ID stored for a source.
func (s *StateStore) GetCursorPosition(source string) (CursorPosition, error) {
	var c Cursor
	err := s.db.First(&c, "source = ?", source).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return CursorPosition{}, nil
		}
		return CursorPosition{}, fmt.Errorf("failed to get cursor: %w", err)
	}
	return CursorPosition{LastTS: c.LastTs, LastID: c.LastID}, nil
}

// SetCursorPosition stores the timestamp and source ID for a source.
func (s *StateStore) SetCursorPosition(source string, position CursorPosition) error {
	c := Cursor{Source: source, LastTs: position.LastTS, LastID: position.LastID}
	err := s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "source"}},
		DoUpdates: clause.AssignmentColumns([]string{"last_ts", "last_id"}),
	}).Create(&c).Error
	if err != nil {
		return fmt.Errorf("failed to set cursor position: %w", err)
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
