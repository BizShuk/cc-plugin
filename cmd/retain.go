package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func retainLogic(store *StateStore, maxAgeDays int, pruneGbrainDir string) error {
	now := time.Now().Unix()
	cutoff := now - int64(maxAgeDays)*86400

	items, err := store.DueForPrune(cutoff)
	if err != nil {
		return err
	}

	prunedCount := 0
	for _, item := range items {
		if item.Source == "gbrain-working" {
			target := filepath.Join(pruneGbrainDir, item.SourceID)
			if _, err := os.Stat(target); err == nil {
				if err := os.Remove(target); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: failed to delete %s: %v\n", target, err)
				} else {
					prunedCount++
				}
			}
		}
		if err := store.DropDistilled(item.Source, item.SourceID); err != nil {
			return fmt.Errorf("failed to drop distilled item %s: %w", item.SourceID, err)
		}
	}

	if prunedCount > 0 {
		fmt.Printf("Successfully pruned %d files and cleaned state distilled entries.\n", prunedCount)
	}
	return nil
}

func RetainCmd() *cobra.Command {
	var statePath string
	var maxAgeDays int
	var pruneGbrainDir string

	cmd := &cobra.Command{
		Use:   "retain",
		Short: "Sweep distilled memories older than max age",
		RunE: func(cmd *cobra.Command, args []string) error {
			if statePath == "" {
				statePath = expandPath(viper.GetString("state.db_path"))
			}
			if maxAgeDays == 0 {
				maxAgeDays = viper.GetInt("retention.max_age_days")
				if maxAgeDays == 0 {
					maxAgeDays = 30
				}
			}
			if pruneGbrainDir == "" {
				pruneGbrainDir = expandPath(viper.GetString("sources.gbrain_working.root"))
			}

			store, err := NewStateStore(statePath)
			if err != nil {
				return err
			}
			defer store.Close()

			return retainLogic(store, maxAgeDays, pruneGbrainDir)
		},
	}

	cmd.Flags().StringVar(&statePath, "state", "", "Path to state.db")
	cmd.Flags().IntVar(&maxAgeDays, "max-age", 0, "Max age in days to retain")
	cmd.Flags().StringVar(&pruneGbrainDir, "prune-gbrain", "", "Path to gbrain/working directory")

	return cmd
}
