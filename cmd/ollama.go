package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bizshuk/cc-plugin/model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ExtractCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extract",
		Short: "Directly extract memories from JSON observations on stdin using Ollama",
		RunE: func(cmd *cobra.Command, args []string) error {
			var observations []model.Observation
			if err := json.NewDecoder(os.Stdin).Decode(&observations); err != nil {
				return fmt.Errorf("failed to parse observations from stdin: %w", err)
			}

			svc := NewOllamaService()
			candidates, err := svc.Extract(observations)
			if err != nil {
				return err
			}

			output, err := json.MarshalIndent(candidates, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(output))
			return nil
		},
	}
	return cmd
}

type OllamaService struct {
	Model   string
	Host    string
	Timeout time.Duration
}

func NewOllamaService() *OllamaService {
	host, model := viper.GetString("llm.host"), viper.GetString("llm.model")
	return &OllamaService{
		Model:   model,
		Host:    strings.TrimSuffix(host, "/"),
		Timeout: 120 * time.Second,
	}
}

const ExtractSystemPrompt = `You extract durable, reusable memories from agent/chat observations. ` +
	`Return ONLY a JSON object: {"candidates": [...]}. Each candidate has: ` +
	`text (verbatim statement worth remembering), ` +
	`entities (list of canonical names: people/projects/topics), ` +
	`kind (one of "fact","experience","preference","inference"), ` +
	`first_person (true if it is the human's own first-person life fact/experience), ` +
	`confirmed_by_human (true only if the human explicitly confirmed it). ` +
	`Omit chit-chat and transient operational noise.`

func (s *OllamaService) Extract(observations []model.Observation) ([]model.Candidate, error) {
	if len(observations) == 0 {
		return nil, nil
	}

	var parts []string
	for _, o := range observations {
		parts = append(parts, fmt.Sprintf("[%s:%s] %s", o.Source, o.SourceID, o.Text))
	}
	joined := strings.Join(parts, "\n\n")

	payload := map[string]interface{}{
		"model": s.Model,
		"messages": []map[string]string{
			{"role": "system", "content": ExtractSystemPrompt},
			{"role": "user", "content": joined},
		},
		"format": "json",
		"stream": false,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	url := s.Host + "/api/chat"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: s.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call Ollama API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Ollama returned status %d: %s", resp.StatusCode, string(body))
	}

	var reply struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&reply); err != nil {
		return nil, fmt.Errorf("failed to decode Ollama reply: %w", err)
	}

	var responseObj struct {
		Candidates []model.Candidate `json:"candidates"`
	}
	if err := json.Unmarshal([]byte(reply.Message.Content), &responseObj); err != nil {
		return nil, fmt.Errorf("failed to parse structured JSON from Ollama content: %w. raw content: %s", err, reply.Message.Content)
	}

	var refs [][]string
	for _, o := range observations {
		refs = append(refs, []string{o.Source, o.SourceID})
	}

	for i := range responseObj.Candidates {
		responseObj.Candidates[i].SourceRefs = refs
	}

	return responseObj.Candidates, nil
}
