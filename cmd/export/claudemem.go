package export

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/bizshuk/cc-plugin/model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	claudeMemSource             = "claude-mem"
	claudeMemExportCursorSource = "claude-mem-export"
)

func claudeMemRead(s *model.StateStore, fromCursor bool) ([]model.Observation, model.CursorPosition, error) {
	lastPosition := model.CursorPosition{}
	if fromCursor {
		var err error
		lastPosition, err = s.GetCursorPosition(claudeMemExportCursorSource)
		if err != nil {
			return nil, model.CursorPosition{}, fmt.Errorf("get claude-mem export cursor: %w", err)
		}
	}

	dbPath := model.ExpandPath(viper.GetString("sources.claude_mem.db_path"))
	dbURL := (&url.URL{Scheme: "file", Path: dbPath, RawQuery: "mode=ro"}).String()
	cmDB, err := gorm.Open(sqlite.Open(dbURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, model.CursorPosition{}, fmt.Errorf("failed to open claude-mem db: %w", err)
	}
	sqlDB, err := cmDB.DB()
	if err != nil {
		return nil, model.CursorPosition{}, fmt.Errorf("failed to access claude-mem db connection: %w", err)
	}

	var dbObs []model.ClaudeMemObservation
	queryErr := cmDB.Where("id > ?", lastPosition.LastID).Order("id ASC").Find(&dbObs).Error
	closeErr := sqlDB.Close()
	if queryErr != nil {
		queryErr = fmt.Errorf("failed to query claude-mem observations: %w", queryErr)
	}
	if closeErr != nil {
		closeErr = fmt.Errorf("failed to close claude-mem db: %w", closeErr)
	}
	if err := errors.Join(queryErr, closeErr); err != nil {
		return nil, model.CursorPosition{}, err
	}

	observations := make([]model.Observation, 0, len(dbObs))
	maxPosition := lastPosition

	for _, o := range dbObs {
		observationID, err := strconv.ParseInt(o.ID, 10, 64)
		if err != nil {
			return nil, model.CursorPosition{}, fmt.Errorf("parse claude-mem observation ID %q: %w", o.ID, err)
		}
		observations = append(observations, model.Observation{
			Source:    claudeMemSource,
			SourceID:  o.ID,
			Timestamp: o.CreatedAtEpoch,
			Text:      o.Text,
		})
		maxPosition = model.CursorPosition{LastTS: o.CreatedAtEpoch, LastID: observationID}
	}

	return observations, maxPosition, nil
}

// ClaudeMemCmd returns the claude-mem export subcommand.
func ClaudeMemCmd() *cobra.Command {
	var allFlag bool

	cmd := &cobra.Command{
		Use:   "claudemem",
		Short: "Export observations from claude-mem SQLite DB",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := model.NewStateStore()
			if err != nil {
				return err
			}
			defer s.Close()

			observations, maxPosition, err := claudeMemRead(s, !allFlag)
			if err != nil {
				return err
			}

			encoder := json.NewEncoder(cmd.OutOrStdout())
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(observations); err != nil {
				return fmt.Errorf("write claude-mem export: %w", err)
			}

			// Update cursor only after export/write finished
			if maxPosition.LastID > 0 {
				lastPosition, err := s.GetCursorPosition(claudeMemExportCursorSource)
				if err != nil {
					return fmt.Errorf("get claude-mem export cursor: %w", err)
				}
				if maxPosition.LastID > lastPosition.LastID {
					if err := s.SetCursorPosition(claudeMemExportCursorSource, maxPosition); err != nil {
						return fmt.Errorf("set claude-mem export cursor: %w", err)
					}
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&allFlag, "all", false, "Export all records from epoch 0 instead of from cursor")

	return cmd
}
