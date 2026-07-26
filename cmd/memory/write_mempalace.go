package memory

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/bizshuk/cc-plugin/model"
	gosdkconfig "github.com/bizshuk/gosdk/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func writeMempalaceLogic(facts []model.Fact, tempDir, wing string) error {
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

	return nil
}

var (
	mempalaceTempDir string
	mempalaceWing    string
)

// WriteMempalaceCmd writes verified facts and runs mempalace mine.
var WriteMempalaceCmd = &cobra.Command{
	Use:   "write-mempalace",
	Short: "Write verified facts to temp files and run mempalace mine",
	RunE: func(cmd *cobra.Command, _ []string) error {
		if mempalaceTempDir == "" {
			mempalaceTempDir = gosdkconfig.ExpandHome(viper.GetString("stores.mempalace.temp_dir"))
		}
		if mempalaceWing == "" {
			mempalaceWing = viper.GetString("stores.mempalace.wing")
		}

		decoder := json.NewDecoder(os.Stdin)
		var facts []model.Fact

		var raw json.RawMessage
		if err := decoder.Decode(&raw); err != nil {
			return fmt.Errorf("failed to decode stdin JSON: %w", err)
		}

		// Try array first
		if err := json.Unmarshal(raw, &facts); err != nil {
			var single model.Fact
			if err2 := json.Unmarshal(raw, &single); err2 != nil {
				return fmt.Errorf("stdin must be a JSON array of Fact or a single Fact object")
			}
			facts = append(facts, single)
		}

		if err := writeMempalaceLogic(facts, mempalaceTempDir, mempalaceWing); err != nil {
			return err
		}

		fmt.Printf("Successfully mined %d facts into mempalace (wing %s).\n", len(facts), mempalaceWing)
		return nil
	},
}

func init() {
	WriteMempalaceCmd.Flags().StringVar(&mempalaceTempDir, "temp-dir", "", "Temporary directory to stage facts")
	WriteMempalaceCmd.Flags().StringVar(&mempalaceWing, "wing", "", "Wing name to mine into")
}
