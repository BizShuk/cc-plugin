package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

func RetainCmd() *cobra.Command {
	var statePath string
	var maxAgeDays int
	var pruneGbrainDir string

	cmd := &cobra.Command{
		Use:   "retain",
		Short: "Sweep distilled memories older than max age",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := NewStateStore(statePath)
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
					return fmt.Errorf("failed to drop distilled item: %w", err)
				}
			}

			fmt.Printf("Successfully pruned %d files and cleaned state distilled entries.\n", prunedCount)
			return nil
		},
	}

	cmd.Flags().StringVar(&statePath, "state", filepath.Join(os.Getenv("HOME"), ".distiller", "state.db"), "Path to state.db")
	cmd.Flags().IntVar(&maxAgeDays, "max-age", 30, "Max age in days to retain")
	cmd.Flags().StringVar(&pruneGbrainDir, "prune-gbrain", filepath.Join(os.Getenv("HOME"), "brain", "working"), "Path to gbrain/working directory")

	return cmd
}
