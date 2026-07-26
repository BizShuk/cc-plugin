package topology

import (
	"fmt"

	topologypkg "github.com/bizshuk/cc-plugin/pkg/topology"
	"github.com/spf13/cobra"
)

const defaultRoot = "plugins/ultra-explore/skills/topology-builder/references"

// TopologyCmd operates on topology-builder knowledge graphs.
var TopologyCmd = &cobra.Command{
	Use:   "topology",
	Short: "Operate on topology-builder knowledge graphs",
}

func init() {
	TopologyCmd.PersistentFlags().String("root", defaultRoot, "topology root directory")
	TopologyCmd.AddCommand(
		VerifyCmd,
		UnlinkedCmd,
		QueryCmd,
		BacklinksCmd,
		IndexCmd,
		RewriteCmd,
	)
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
