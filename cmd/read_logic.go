package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bizshuk/cc-plugin/model"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func readGbrainLogic(store *StateStore, workingDir string) ([]model.Observation, int64, error) {
	lastTS, err := store.GetCursor("gbrain-working")
	if err != nil {
		return nil, 0, err
	}

	var observations []model.Observation
	var maxTS = lastTS

	err = filepath.Walk(workingDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			mtime := info.ModTime().Unix()
			if mtime > lastTS {
				content, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				rel, err := filepath.Rel(workingDir, path)
				if err != nil {
					return err
				}
				observations = append(observations, model.Observation{
					Source:    "gbrain-working",
					SourceID:  rel,
					Timestamp: mtime,
					Text:      string(content),
				})
				if mtime > maxTS {
					maxTS = mtime
				}
			}
		}
		return nil
	})
	if err != nil && !os.IsNotExist(err) {
		return nil, 0, err
	}

	return observations, maxTS, nil
}

func readClaudeMemLogic() ([]model.Observation, int64, error) {
	store, err := NewStateStore()
	if err != nil {
		return nil, 0, err
	}
	defer store.Close()

	lastTS, err := store.GetCursor("claude-mem")
	if err != nil {
		return nil, 0, err
	}

	dbPath := expandPath(viper.GetString("sources.claude_mem.db_path"))
	cmDB, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to open claude-mem db: %w", err)
	}

	var dbObs []model.ClaudeMemObservation
	err = cmDB.Where("created_at_epoch > ?", lastTS).Order("created_at_epoch ASC").Find(&dbObs).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query claude-mem observations: %w", err)
	}

	var observations []model.Observation
	maxTS := lastTS

	for _, o := range dbObs {
		observations = append(observations, model.Observation{
			Source:    "claude-mem",
			SourceID:  o.ID,
			Timestamp: o.CreatedAtEpoch,
			Text:      o.Text,
		})

		if o.CreatedAtEpoch > maxTS {
			maxTS = o.CreatedAtEpoch
		}
	}

	return observations, maxTS, nil
}
