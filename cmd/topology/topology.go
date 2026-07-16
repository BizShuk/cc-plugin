package topology

import (
	"fmt"

	topologypkg "github.com/bizshuk/cc-plugin/pkg/topology"
	"github.com/spf13/cobra"
)

const defaultRoot = "plugins/ultra-explore/skills/topology-builder/references"

// TopologyCmd returns the top-level topology command.
func TopologyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "topology",
		Short: "Operate on topology-builder knowledge graphs",
	}
	cmd.PersistentFlags().String("root", defaultRoot, "topology root directory")
	cmd.AddCommand(
		VerifyCmd(),
		UnlinkedCmd(),
		QueryCmd(),
		BacklinksCmd(),
		IndexCmd(),
		RewriteCmd(),
	)
	return cmd
}

func loadFromFlags(cmd *cobra.Command) (*topologypkg.Topology, error) {
	root, err := cmd.Flags().GetString("root")
	if err != nil {
		return nil, fmt.Errorf("read root flag: %w", err)
	}
	topo, err := topologypkg.LoadTopology(root)
	if err != nil {
		return nil, fmt.Errorf("load topology: %w", err)
	}
	return topo, nil
}
