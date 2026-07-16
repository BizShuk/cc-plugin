package cmd

import (
	"log/slog"
	"os"

	"github.com/bizshuk/cc-plugin/cmd/export"
	"github.com/bizshuk/cc-plugin/cmd/memory"
	"github.com/bizshuk/cc-plugin/cmd/topology"
	"github.com/bizshuk/cc-plugin/config"
	"github.com/spf13/cobra"
)

var RootCmd *cobra.Command

func init() {
	config.Init()

	RootCmd = &cobra.Command{
		Use:           "cc-plugin",
		Short:         "Manage memory distillation, exports, and plugin utilities",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	RootCmd.AddCommand(memory.DistillCmd())
	RootCmd.AddCommand(memory.RetainCmd())
	RootCmd.AddCommand(memory.WriteAgentMemoryCmd())
	RootCmd.AddCommand(memory.WriteMempalaceCmd())
	RootCmd.AddCommand(memory.ExtractCmd())
	RootCmd.AddCommand(memory.ResetCmd())
	RootCmd.AddCommand(export.ExportCmd())
	RootCmd.AddCommand(topology.TopologyCmd())
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		slog.Error("command failed", "err", err)
		os.Exit(1)
	}
}
