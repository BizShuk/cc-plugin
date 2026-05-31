package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func QualifiesForTruth(c Candidate, corroboration int) bool {
	if c.Kind == "inference" {
		return false
	}
	if c.ConfirmedByHuman {
		return true
	}
	if c.FirstPerson && (c.Kind == "fact" || c.Kind == "experience") {
		return true
	}
	if corroboration >= 2 {
		return true
	}
	return false
}

func DistillCmd() *cobra.Command {
	var noRetain bool

	cmd := &cobra.Command{
		Use:   "distill",
		Short: "Distill memories from source databases and local markdown notes",
		RunE: func(cmd *cobra.Command, args []string) error {
			statePath := expandPath(viper.GetString("state.db_path"))
			store, err := NewStateStore(statePath)
			if err != nil {
				return err
			}
			defer store.Close()

			// 1. Read gbrain observations
			gbrainRoot := expandPath(viper.GetString("sources.gbrain_working.root"))
			gbrainObs, gbrainMaxTS, err := readGbrainLogic(store, gbrainRoot)
			if err != nil {
				return err
			}

			// 2. Read claude-mem observations
			cmDB := expandPath(viper.GetString("sources.claude_mem.db_path"))
			cmTable := viper.GetString("sources.claude_mem.table")
			cmIdCol := viper.GetString("sources.claude_mem.id_col")
			cmTsCol := viper.GetString("sources.claude_mem.ts_col")
			cmTextCol := viper.GetString("sources.claude_mem.text_col")
			cmObs, cmMaxTS, err := readClaudeMemLogic(store, cmDB, cmTable, cmIdCol, cmTsCol, cmTextCol)
			if err != nil {
				return err
			}

			allObs := append(gbrainObs, cmObs...)
			if len(allObs) == 0 {
				fmt.Println("[distiller] No new observations found.")
				return nil
			}

			// 3. Extract candidates via Ollama Service
			llm := NewOllamaService()
			candidates, err := llm.Extract(allObs)
			if err != nil {
				return err
			}

			// 4. Process candidates into Memories and Facts
			var memories []Memory
			var facts []Fact
			now := time.Now().Unix()

			for _, c := range candidates {
				fp := Fingerprint(c.Text, c.Entities)

				var corroboration int
				for _, ref := range c.SourceRefs {
					if len(ref) > 0 {
						count, err := store.RecordSeen(fp, ref[0])
						if err != nil {
							return err
						}
						corroboration = count
					}
				}

				memories = append(memories, Memory{
					Fingerprint: fp,
					Text:        c.Text,
					Entities:    c.Entities,
					Kind:        c.Kind,
					CreatedAt:   now,
				})

				if QualifiesForTruth(c, corroboration) {
					facts = append(facts, Fact{
						Fingerprint: fp,
						Text:        c.Text,
						Entities:    c.Entities,
						Evidence:    c.SourceRefs,
						CreatedAt:   now,
					})
				}
			}

			// Write to agentmemory API
			if len(memories) > 0 {
				url := viper.GetString("stores.agentmemory.url")
				if err := writeAgentMemoryLogic(memories, url); err != nil {
					return err
				}
			}

			// Write to mempalace CLI
			if len(facts) > 0 {
				tempDir := expandPath(viper.GetString("stores.mempalace.temp_dir"))
				wing := viper.GetString("stores.mempalace.wing")
				if err := writeMempalaceLogic(facts, tempDir, wing); err != nil {
					return err
				}
			}

			// Mark observations as distilled
			for _, o := range allObs {
				if err := store.MarkDistilled(o.Source, o.SourceID, now); err != nil {
					return err
				}
			}

			// Update cursors
			if gbrainMaxTS > 0 {
				if err := store.SetCursor("gbrain-working", gbrainMaxTS); err != nil {
					return err
				}
			}
			if cmMaxTS > 0 {
				if err := store.SetCursor("claude-mem", cmMaxTS); err != nil {
					return err
				}
			}

			fmt.Printf("[distiller] Pipeline ran successfully. Sources read=%d, Memories written=%d, Facts written=%d\n", len(allObs), len(memories), len(facts))

			// 5. Sweep retention
			if !noRetain {
				maxAgeDays := viper.GetInt("retention.max_age_days")
				if maxAgeDays == 0 {
					maxAgeDays = 30
				}
				pruneGbrainDir := expandPath(viper.GetString("sources.gbrain_working.root"))
				if err := retainLogic(store, maxAgeDays, pruneGbrainDir); err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&noRetain, "no-retain", false, "Disable retention sweep after run")
	return cmd
}
