package cmd

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

func ReadClaudeMemCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read-claudemem",
		Short: "Read observations from claude-mem SQLite DB",
		RunE: func(cmd *cobra.Command, args []string) error {
			observations, maxTS, err := readClaudeMemLogic()
			if err != nil {
				return err
			}

			output, err := json.MarshalIndent(observations, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(output))

			if maxTS > 0 {
				store, err := NewStateStore()
				if err != nil {
					return err
				}
				defer store.Close()

				lastCursor, _ := store.GetCursor("claude-mem")
				if maxTS > lastCursor {
					if err := store.SetCursor("claude-mem", maxTS); err != nil {
						return err
					}
				}
			}

			return nil
		},
	}

	return cmd
}
