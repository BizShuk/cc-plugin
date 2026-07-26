package export

import (
	"github.com/spf13/cobra"
)

// ExportCmd exports data from the configured source commands.
var ExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export data from various sources",
}

func init() {
	ExportCmd.AddCommand(ClaudeMemCmd)
}
