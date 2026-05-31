package export

import (
	"github.com/spf13/cobra"
)

// ExportCmd returns the top-level export Cobra command.
func ExportCmd() *cobra.Command {
	exportCmd := &cobra.Command{
		Use:   "export",
		Short: "Export data from various sources",
	}

	exportCmd.AddCommand(ClaudeMemCmd())
	exportCmd.AddCommand(GbrainCmd())
	exportCmd.AddCommand(MempalaceCmd())

	return exportCmd
}
