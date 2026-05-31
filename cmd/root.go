package cmd

import (
	"fmt"
	"os"

	"github.com/bizshuk/cc-plugin/config"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "distiller",
	Short: "Distiller CLI manages memory systems cross laptop and server",
}

func init() {
	config.Init()

	RootCmd.AddCommand(RetainCmd())
	RootCmd.AddCommand(ReadGbrainCmd())
	RootCmd.AddCommand(ReadClaudeMemCmd())
	RootCmd.AddCommand(WriteAgentMemoryCmd())
	RootCmd.AddCommand(WriteMempalaceCmd())
	RootCmd.AddCommand(ExtractCmd())
	RootCmd.AddCommand(DistillCmd())
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
