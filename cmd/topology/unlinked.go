package topology

import (
	"fmt"

	"github.com/spf13/cobra"
)

// UnlinkedCmd returns the unlinked-entity report command.
func UnlinkedCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unlinked",
		Short: "List entities with no inbound or outbound edges",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			topo, err := loadFromFlags(cmd)
			if err != nil {
				return err
			}
			noInbound, noOutbound := topo.Unlinked()
			printEntityList(cmd, "no inbound", noInbound)
			printEntityList(cmd, "no outbound", noOutbound)
			return nil
		},
	}
}

func printEntityList(cmd *cobra.Command, label string, names []string) {
	fmt.Fprintf(cmd.OutOrStdout(), "%s:\n", label)
	if len(names) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "  (none)")
		return
	}
	for _, name := range names {
		fmt.Fprintf(cmd.OutOrStdout(), "  %s\n", name)
	}
}
