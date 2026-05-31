package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func writeAgentMemoryLogic(memories []Memory, url string) error {
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
	return nil
}

func WriteAgentMemoryCmd() *cobra.Command {
	var url string

	cmd := &cobra.Command{
		Use:   "write-agentmemory",
		Short: "Post distilled memories from stdin into agentmemory API",
		RunE: func(cmd *cobra.Command, args []string) error {
			if url == "" {
				url = viper.GetString("stores.agentmemory.url")
			}

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

			if err := writeAgentMemoryLogic(memories, url); err != nil {
				return err
			}

			fmt.Printf("Successfully wrote %d memories into agentmemory.\n", len(memories))
			return nil
		},
	}

	cmd.Flags().StringVar(&url, "url", "", "agentmemory API remember endpoint")

	return cmd
}
