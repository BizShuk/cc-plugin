package memory

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/bizshuk/cc-plugin/model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RetainCmd() *cobra.Command {
	var maxAgeDays int
	var pruneGbrainDir string

	cmd := &cobra.Command{
		Use:   "retain",
		Short: "Sweep distilled memories older than max age",
		RunE: func(cmd *cobra.Command, args []string) error {
			return retainLogic()
		},
	}

	cmd.Flags().IntVar(&maxAgeDays, "max-age", 0, "Max age in days to retain")
	cmd.Flags().StringVar(&pruneGbrainDir, "prune-gbrain", "", "Path to gbrain/working directory")

	return cmd
}

func retainLogic() error {
	maxAgeDays := viper.GetInt("retention.max_age_days")
	pruneGbrainDir := model.ExpandPath(viper.GetString("sources.gbrain_working.root"))

	store, err := NewStateStore()
	if err != nil {
		return err
	}
	defer store.Close()

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
			fmt.Fprintf(os.Stderr, "Warning: failed to drop distilled item %s: %v\n", item.SourceID, err)
		}
	}

	if prunedCount > 0 {
		fmt.Printf("Successfully pruned %d files and cleaned state distilled entries.\n", prunedCount)
	}
	return nil
}
