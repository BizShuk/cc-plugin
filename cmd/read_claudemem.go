package cmd

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func readClaudeMemLogic(store *StateStore, dbPath, table, idCol, tsCol, textCol string) ([]Observation, int64, error) {
	lastTS, err := store.GetCursor("claude-mem")
	if err != nil {
		return nil, 0, err
	}

	cmDB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to open claude-mem db: %w", err)
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
			return nil, 0, fmt.Errorf("failed to query claude-mem observations: %w", err)
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
			return nil, 0, fmt.Errorf("failed to scan observation: %w", err)
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

	return observations, maxTS, nil
}

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
			if statePath == "" {
				statePath = expandPath(viper.GetString("state.db_path"))
			}
			if dbPath == "" {
				dbPath = expandPath(viper.GetString("sources.claude_mem.db_path"))
			}
			if table == "" {
				table = viper.GetString("sources.claude_mem.table")
			}
			if idCol == "" {
				idCol = viper.GetString("sources.claude_mem.id_col")
			}
			if tsCol == "" {
				tsCol = viper.GetString("sources.claude_mem.ts_col")
			}
			if textCol == "" {
				textCol = viper.GetString("sources.claude_mem.text_col")
			}

			store, err := NewStateStore(statePath)
			if err != nil {
				return err
			}
			defer store.Close()

			observations, maxTS, err := readClaudeMemLogic(store, dbPath, table, idCol, tsCol, textCol)
			if err != nil {
				return err
			}

			output, err := json.MarshalIndent(observations, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(output))

			if maxTS > 0 {
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

	cmd.Flags().StringVar(&statePath, "state", "", "Path to state.db")
	cmd.Flags().StringVar(&dbPath, "db", "", "Path to claude-mem database file")
	cmd.Flags().StringVar(&table, "table", "", "Observations table name")
	cmd.Flags().StringVar(&idCol, "id-col", "", "ID column name")
	cmd.Flags().StringVar(&tsCol, "ts-col", "", "Timestamp column name")
	cmd.Flags().StringVar(&textCol, "text-col", "", "Text content column name")

	return cmd
}
