package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

func WriteMempalaceCmd() *cobra.Command {
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
