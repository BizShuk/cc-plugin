package memory

import (
	"fmt"
	"time"

	"github.com/bizshuk/cc-plugin/model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func QualifiesForTruth(c model.Candidate, corroboration int) bool {
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
			store, err := NewStateStore()
			if err != nil {
				return err
			}
			defer store.Close()

			// 1. Read gbrain observations
			gbrainRoot := model.ExpandPath(viper.GetString("sources.gbrain_working.root"))
			gbrainObs, gbrainMaxTS, err := readGbrainLogic(store, gbrainRoot)
			if err != nil {
				return err
			}

			// 2. Read claude-mem observations
			cmObs, cmMaxTS, err := readClaudeMemLogic()
			if err != nil {
				return err
			}

			allObs := append(gbrainObs, cmObs...)
			if len(allObs) == 0 {
				fmt.Println("No new observations found.")
				return nil
			}

			// 3. Extract candidates via Ollama Service
			llm := NewOllamaService()
			candidates, err := llm.Extract(cmd.Context(), allObs)
			if err != nil {
				return err
			}

			// 4. Process candidates into Memories and Facts
			var memories []model.Memory
			var facts []model.Fact
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

				memories = append(memories, model.Memory{
					Fingerprint: fp,
					Text:        c.Text,
					Entities:    c.Entities,
					Kind:        c.Kind,
					CreatedAt:   now,
				})

				if QualifiesForTruth(c, corroboration) {
					facts = append(facts, model.Fact{
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
				tempDir := model.ExpandPath(viper.GetString("stores.mempalace.temp_dir"))
				wing := viper.GetString("stores.mempalace.wing")
				if err := writeMempalaceLogic(facts, tempDir, wing); err != nil {
					return err
				}
			}

			// Mark observations as distilled
			var distilledItems []model.DistilledItem
			for _, o := range allObs {
				distilledItems = append(distilledItems, model.DistilledItem{Source: o.Source, SourceID: o.SourceID})
			}
			if err := store.MarkDistilledBatch(distilledItems, now); err != nil {
				return err
			}

			// Update cursors
			cursors := make(map[string]int64)
			if gbrainMaxTS > 0 {
				cursors["gbrain-working"] = gbrainMaxTS
			}
			if cmMaxTS > 0 {
				cursors["claude-mem"] = cmMaxTS
			}
			if len(cursors) > 0 {
				if err := store.SetCursorsBatch(cursors); err != nil {
					return err
				}
			}

			fmt.Printf("[distiller] Pipeline ran successfully. Sources read=%d, Memories written=%d, Facts written=%d\n", len(allObs), len(memories), len(facts))

			// 5. Sweep retention
			if !noRetain {
				if err := retainLogic(); err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&noRetain, "no-retain", false, "Disable retention sweep after run")
	return cmd
}
