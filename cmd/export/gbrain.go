package export

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

func gbrainRead(s *model.StateStore, workingDir string, fromCursor bool) ([]model.Observation, int64, error) {
	lastTS := int64(0)
	if fromCursor {
		var err error
		lastTS, err = s.GetCursor("gbrain-working")
		if err != nil {
			return nil, 0, err
		}
	}

	var observations []model.Observation
	var maxTS = lastTS

	err := filepath.Walk(workingDir, func(path string, info os.FileInfo, err error) error {
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

// GbrainCmd returns the gbrain export subcommand.
func GbrainCmd() *cobra.Command {
	var dir string

	cmd := &cobra.Command{
		Use:   "gbrain",
		Short: "Export markdown logs from gbrain/working directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			if dir == "" {
				dir = model.ExpandPath(viper.GetString("sources.gbrain_working.root"))
			}

			s, err := model.NewStateStore()
			if err != nil {
				return err
			}
			defer s.Close()

			observations, maxTS, err := gbrainRead(s, dir, !cmd.Flags().Changed("all"))
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
				lastCursor, _ := s.GetCursor("gbrain-working")
				if maxTS > lastCursor {
					if err := s.SetCursor("gbrain-working", maxTS); err != nil {
						return err
					}
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&dir, "dir", "", "Path to gbrain/working directory")
	cmd.Flags().Bool("all", false, "Export all records from epoch 0 instead of from cursor (stored in lastTS)")

	return cmd
}
