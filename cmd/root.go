package cmd

import (
	"fmt"
	"os"

	"github.com/bizshuk/gosdk/config"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "distiller",
	Short: "Distiller CLI manages memory systems cross laptop and server",
}

const defaultSettings = `{
  "state": {
    "db_path": "~/.distiller/state.db"
  },
  "retention": {
    "max_age_days": 30
  },
  "llm": {
    "model": "qwen2.5",
    "host": "http://localhost:11434"
  },
  "sources": {
    "claude_mem": {
      "db_path": "~/.claude-mem/claude-mem.db",
      "table": "observations",
      "id_col": "id",
      "ts_col": "created_at_epoch",
      "text_col": "text"
    },
    "gbrain_working": {
      "root": "~/brain/working"
    }
  },
  "stores": {
    "agentmemory": {
      "url": "http://localhost:3111/agentmemory/remember"
    },
    "mempalace": {
      "wing": "main",
      "temp_dir": "/tmp/mempalace-temp"
    }
  }
}`

func init() {
	config.Default(
		config.WithAppName("cc-plugin"),
		config.WithDefaultValue(defaultSettings),
	)

	RootCmd.AddCommand(RetainCmd())
	RootCmd.AddCommand(ReadGbrainCmd())
	RootCmd.AddCommand(ReadClaudeMemCmd())
	RootCmd.AddCommand(WriteAgentMemoryCmd())
	RootCmd.AddCommand(WriteMempalaceCmd())
	RootCmd.AddCommand(ExtractCmd())
	RootCmd.AddCommand(RunCmd())
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
