package export

import (
	"encoding/json"
	"fmt"

	"github.com/bizshuk/cc-plugin/model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func claudeMemRead(s *model.StateStore, fromCursor bool) ([]model.Observation, int64, error) {
	lastTS := int64(0)
	if fromCursor {
		var err error
		lastTS, err = s.GetCursor("claude-mem")
		if err != nil {
			return nil, 0, err
		}
	}

	dbPath := model.ExpandPath(viper.GetString("sources.claude_mem.db_path"))
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

			observations, maxTS, err := claudeMemRead(s, !allFlag)
			if err != nil {
				return err
			}

			output, err := json.MarshalIndent(observations, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(output))

			// Update cursor only after export/write finished
			if maxTS > 0 {
				lastCursor, _ := s.GetCursor("claude-mem")
				if maxTS > lastCursor {
					if err := s.SetCursor("claude-mem", maxTS); err != nil {
						return err
					}
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&allFlag, "all", false, "Export all records from epoch 0 instead of from cursor")

	return cmd
}
