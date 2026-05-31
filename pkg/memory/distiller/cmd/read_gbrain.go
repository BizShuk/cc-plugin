package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func ReadGbrainCmd() *cobra.Command {
	var statePath string
	var workingDir string

	cmd := &cobra.Command{
		Use:   "read-gbrain",
		Short: "Read new markdown logs from gbrain/working and update cursor",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := NewStateStore(statePath)
			if err != nil {
				return err
			}
			defer store.Close()

			lastTS, err := store.GetCursor("gbrain-working")
			if err != nil {
				return err
			}

			var observations []Observation
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
						observations = append(observations, Observation{
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
				return err
			}

			// Output observations as JSON
			output, err := json.MarshalIndent(observations, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(output))

			// Update cursor if we read new items
			if maxTS > lastTS {
				if err := store.SetCursor("gbrain-working", maxTS); err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&statePath, "state", filepath.Join(os.Getenv("HOME"), ".distiller", "state.db"), "Path to state.db")
	cmd.Flags().StringVar(&workingDir, "dir", filepath.Join(os.Getenv("HOME"), "brain", "working"), "Path to gbrain/working directory")

	return cmd
}
