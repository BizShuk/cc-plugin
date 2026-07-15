package cmd

import (
	"fmt"
	"os"

	"github.com/bizshuk/cc-plugin/cmd/export"
	"github.com/bizshuk/cc-plugin/cmd/memory"
	"github.com/bizshuk/cc-plugin/config"
	"github.com/spf13/cobra"
)

var RootCmd *cobra.Command

func init() {
	config.Init()

	RootCmd = memory.DistillCmd()

	RootCmd.AddCommand(memory.RetainCmd())
	RootCmd.AddCommand(memory.WriteAgentMemoryCmd())
	RootCmd.AddCommand(memory.WriteMempalaceCmd())
	RootCmd.AddCommand(memory.ExtractCmd())
	RootCmd.AddCommand(memory.ResetCmd())
	RootCmd.AddCommand(export.ExportCmd())
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
