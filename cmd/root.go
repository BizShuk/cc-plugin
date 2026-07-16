package cmd

import (
	"log/slog"
	"os"

	"github.com/bizshuk/cc-plugin/cmd/export"
	"github.com/bizshuk/cc-plugin/cmd/topology"
	"github.com/bizshuk/cc-plugin/config"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "cc-plugin",
	Short: "cc-plugin is a CLI tool for AI agent memory integration and plugins",
	Long:  `A modular CLI tool designed to distill observations, manage plugin skills, and export memories.`,
}

func init() {
	config.Init()

	RootCmd.AddCommand(export.ExportCmd())
	RootCmd.AddCommand(topology.TopologyCmd())
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		slog.Error("command failed", "err", err)
		os.Exit(1)
	}
}
