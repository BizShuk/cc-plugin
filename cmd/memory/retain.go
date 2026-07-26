package memory

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	gosdkconfig "github.com/bizshuk/gosdk/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	retainMaxAgeDays int
	retainGbrainDir  string
)

// RetainCmd sweeps distilled memories older than the configured maximum age.
var RetainCmd = &cobra.Command{
	Use:   "retain",
	Short: "Sweep distilled memories older than max age",
	RunE: func(cmd *cobra.Command, _ []string) error {
		return retainLogic()
	},
}

func init() {
	RetainCmd.Flags().IntVar(&retainMaxAgeDays, "max-age", 0, "Max age in days to retain")
	RetainCmd.Flags().StringVar(&retainGbrainDir, "prune-gbrain", "", "Path to gbrain/working directory")
}

func retainLogic() error {
	maxAgeDays := viper.GetInt("retention.max_age_days")
	pruneGbrainDir := gosdkconfig.ExpandHome(viper.GetString("sources.gbrain_working.root"))

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
					slog.Warn("failed to delete distilled source", "path", target, "err", err)
				} else {
					prunedCount++
				}
			}
		}
		if err := store.DropDistilled(item.Source, item.SourceID); err != nil {
			slog.Warn("failed to drop distilled item", "source", item.Source, "source_id", item.SourceID, "err", err)
		}
	}

	if prunedCount > 0 {
		fmt.Printf("Successfully pruned %d files and cleaned state distilled entries.\n", prunedCount)
	}
	return nil
}
