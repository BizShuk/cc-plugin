package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type Observation struct {
	Source    string `json:"source"`
	SourceID  string `json:"source_id"`
	Timestamp int64  `json:"timestamp"`
	Text      string `json:"text"`
}

type Memory struct {
	Fingerprint string   `json:"fingerprint"`
	Text        string   `json:"text"`
	Entities    []string `json:"entities"`
	Kind        string   `json:"kind"`
	CreatedAt   int64    `json:"created_at"`
}

type Fact struct {
	Fingerprint string     `json:"fingerprint"`
	Text        string     `json:"text"`
	Entities    []string   `json:"entities"`
	Evidence    [][]string `json:"evidence"`
	CreatedAt   int64      `json:"created_at"`
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "distiller",
		Short: "Distiller CLI manages memory systems cross laptop and server",
	}

	rootCmd.AddCommand(retainCmd())
	rootCmd.AddCommand(readGbrainCmd())
	rootCmd.AddCommand(readClaudeMemCmd())
	rootCmd.AddCommand(writeAgentMemoryCmd())
	rootCmd.AddCommand(writeMempalaceCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// 1. Phase 7: retainCmd
func retainCmd() *cobra.Command {
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

// 2. Phase 8: readGbrainCmd
func readGbrainCmd() *cobra.Command {
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

// 3. Phase 9: readClaudeMemCmd
func readClaudeMemCmd() *cobra.Command {
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

// 4. Phase 10 (A): writeAgentMemoryCmd
func writeAgentMemoryCmd() *cobra.Command {
	var url string

	cmd := &cobra.Command{
		Use:   "write-agentmemory",
		Short: "Post distilled memories from stdin into agentmemory API",
		RunE: func(cmd *cobra.Command, args []string) error {
			decoder := json.NewDecoder(os.Stdin)
			var memories []Memory

			// Read stdin as a JSON array or single object
			var raw json.RawMessage
			if err := decoder.Decode(&raw); err != nil {
				return fmt.Errorf("failed to decode stdin JSON: %w", err)
			}

			// Try array first
			if err := json.Unmarshal(raw, &memories); err != nil {
				var single Memory
				if err2 := json.Unmarshal(raw, &single); err2 != nil {
					return fmt.Errorf("stdin must be a JSON array of Memory or a single Memory object")
				}
				memories = append(memories, single)
			}

			for _, mem := range memories {
				// Map to agentmemory format
				payload := map[string]interface{}{
					"content":  mem.Text,
					"concepts": append([]string{mem.Kind}, mem.Entities...),
					"files":    []string{},
				}

				payloadBytes, err := json.Marshal(payload)
				if err != nil {
					return err
				}

				resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
				if err != nil {
					return fmt.Errorf("failed to send remember request to agentmemory: %w", err)
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
					body, _ := io.ReadAll(resp.Body)
					return fmt.Errorf("agentmemory server returned error status %d: %s", resp.StatusCode, string(body))
				}
			}

			fmt.Printf("Successfully wrote %d memories into agentmemory.\n", len(memories))
			return nil
		},
	}

	cmd.Flags().StringVar(&url, "url", "http://localhost:3111/agentmemory/remember", "agentmemory API remember endpoint")

	return cmd
}

// 5. Phase 10 (B): writeMempalaceCmd
func writeMempalaceCmd() *cobra.Command {
	var tempDir string
	var wing string

	cmd := &cobra.Command{
		Use:   "write-mempalace",
		Short: "Write verified facts to temp files and run mempalace mine",
		RunE: func(cmd *cobra.Command, args []string) error {
			decoder := json.NewDecoder(os.Stdin)
			var facts []Fact

			var raw json.RawMessage
			if err := decoder.Decode(&raw); err != nil {
				return fmt.Errorf("failed to decode stdin JSON: %w", err)
			}

			// Try array first
			if err := json.Unmarshal(raw, &facts); err != nil {
				var single Fact
				if err2 := json.Unmarshal(raw, &single); err2 != nil {
					return fmt.Errorf("stdin must be a JSON array of Fact or a single Fact object")
				}
				facts = append(facts, single)
			}

			// Create temp room directory
			roomDir := filepath.Join(tempDir, "general")
			if err := os.MkdirAll(roomDir, 0755); err != nil {
				return err
			}

			for _, fact := range facts {
				evidenceStr, _ := json.Marshal(fact.Evidence)
				entitiesStr, _ := json.Marshal(fact.Entities)
				content := fmt.Sprintf("# Fact: %s\n\n%s\n\nEntities: %s\nEvidence: %s\n", fact.Fingerprint, fact.Text, string(entitiesStr), string(evidenceStr))

				filePath := filepath.Join(roomDir, fmt.Sprintf("%s.md", fact.Fingerprint))
				if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
					return fmt.Errorf("failed to write fact file: %w", err)
				}
			}

			// Initialize mempalace if mempalace.yaml doesn't exist
			yamlPath := filepath.Join(tempDir, "mempalace.yaml")
			if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
				// Run in non-interactive if possible, or create default config
				// We can just create a basic mempalace.yaml ourselves to make it fully automated and non-interactive!
				defaultYaml := fmt.Sprintf("wing: %s\nrooms:\n  general: ['*.md']\n", wing)
				if err := os.WriteFile(yamlPath, []byte(defaultYaml), 0644); err != nil {
					return fmt.Errorf("failed to write mempalace.yaml: %w", err)
				}
			}

			// Run mempalace mine
			mineCmd := exec.Command("mempalace", "mine", tempDir, "--wing", wing)
			var stdout, stderr bytes.Buffer
			mineCmd.Stdout = &stdout
			mineCmd.Stderr = &stderr
			if err := mineCmd.Run(); err != nil {
				return fmt.Errorf("mempalace mine failed: %w\nstdout: %s\nstderr: %s", err, stdout.String(), stderr.String())
			}

			fmt.Printf("Successfully mined %d facts into mempalace (wing %s).\n", len(facts), wing)
			return nil
		},
	}

	cmd.Flags().StringVar(&tempDir, "temp-dir", "/tmp/mempalace-temp", "Temporary directory to stage facts")
	cmd.Flags().StringVar(&wing, "wing", "main", "Wing name to mine into")

	return cmd
}
