package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bizshuk/cc-plugin/model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

func ReadGbrainCmd() *cobra.Command {
	var workingDir string

	cmd := &cobra.Command{
		Use:   "read-gbrain",
		Short: "Read new markdown logs from gbrain/working and update cursor",
		RunE: func(cmd *cobra.Command, args []string) error {
			if workingDir == "" {
				workingDir = expandPath(viper.GetString("sources.gbrain_working.root"))
			}

			store, err := NewStateStore()
			if err != nil {
				return err
			}
			defer store.Close()

			observations, maxTS, err := readGbrainLogic(store, workingDir)
			if err != nil {
				return err
			}

			// Output observations as JSON
			output, err := json.MarshalIndent(observations, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(output))

			// Update cursor if we read new items
			if maxTS > viper.GetInt64("state.gbrain_working.cursor") && maxTS > 0 {
				lastTS, _ := store.GetCursor("gbrain-working")
				if maxTS > lastTS {
					if err := store.SetCursor("gbrain-working", maxTS); err != nil {
						return err
					}
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&workingDir, "dir", "", "Path to gbrain/working directory")

	return cmd
}
