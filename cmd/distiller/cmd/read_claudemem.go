package cmd

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func ReadClaudeMemCmd() *cobra.Command {
	var statePath string
	var dbPath string
	var table string
	var idCol string
	var tsCol string
	var textCol string

	cmd := &cobra.Command{
		Use:   "read-claudemem",
		Short: "Read observations from claude-mem SQLite DB",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := NewStateStore(statePath)
			if err != nil {
				return err
			}
			defer store.Close()

			lastTS, err := store.GetCursor("claude-mem")
			if err != nil {
				return err
			}

			cmDB, err := sql.Open("sqlite3", dbPath)
			if err != nil {
				return fmt.Errorf("failed to open claude-mem db: %w", err)
			}
			defer cmDB.Close()

			// Check if columns exist to construct detailed text
			query := fmt.Sprintf(
				"SELECT %s, %s, %s, title, subtitle, facts, narrative FROM %s WHERE %s > ? ORDER BY %s ASC",
				idCol, tsCol, textCol, table, tsCol, tsCol,
			)

			rows, err := cmDB.Query(query, lastTS)
			if err != nil {
				// Fallback to simpler query if extra columns don't exist
				query = fmt.Sprintf(
					"SELECT %s, %s, %s FROM %s WHERE %s > ? ORDER BY %s ASC",
					idCol, tsCol, textCol, table, tsCol, tsCol,
				)
				rows, err = cmDB.Query(query, lastTS)
				if err != nil {
					return fmt.Errorf("failed to query claude-mem observations: %w", err)
				}
			}
			defer rows.Close()

			var observations []Observation
			var maxTS = lastTS

			for rows.Next() {
				var sid string
				var ts int64
				var textVal sql.NullString
				var title, subtitle, facts, narrative sql.NullString

				cols, _ := rows.Columns()
				if len(cols) > 3 {
					err = rows.Scan(&sid, &ts, &textVal, &title, &subtitle, &facts, &narrative)
				} else {
					err = rows.Scan(&sid, &ts, &textVal)
				}
				if err != nil {
					return fmt.Errorf("failed to scan observation: %w", err)
				}

				var fullText string
				if textVal.Valid && textVal.String != "" {
					fullText = textVal.String
				} else {
					var parts []string
					if title.Valid && title.String != "" {
						parts = append(parts, "Title: "+title.String)
					}
					if subtitle.Valid && subtitle.String != "" {
						parts = append(parts, "Subtitle: "+subtitle.String)
					}
					if narrative.Valid && narrative.String != "" {
						parts = append(parts, "Narrative: "+narrative.String)
					}
					if facts.Valid && facts.String != "" {
						parts = append(parts, "Facts: "+facts.String)
					}
					fullText = strings.Join(parts, "\n")
				}

				observations = append(observations, Observation{
					Source:    "claude-mem",
					SourceID:  sid,
					Timestamp: ts,
					Text:      fullText,
				})

				if ts > maxTS {
					maxTS = ts
				}
			}

			output, err := json.MarshalIndent(observations, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(output))

			if maxTS > lastTS {
				if err := store.SetCursor("claude-mem", maxTS); err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&statePath, "state", filepath.Join(os.Getenv("HOME"), ".distiller", "state.db"), "Path to state.db")
	cmd.Flags().StringVar(&dbPath, "db", filepath.Join(os.Getenv("HOME"), ".claude-mem", "claude-mem.db"), "Path to claude-mem database file")
	cmd.Flags().StringVar(&table, "table", "observations", "Observations table name")
	cmd.Flags().StringVar(&idCol, "id-col", "id", "ID column name")
	cmd.Flags().StringVar(&tsCol, "ts-col", "created_at_epoch", "Timestamp column name")
	cmd.Flags().StringVar(&textCol, "text-col", "text", "Text content column name")

	return cmd
}
