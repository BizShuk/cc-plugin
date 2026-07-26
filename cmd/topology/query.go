package topology

import (
	"fmt"
	"sort"

	topologypkg "github.com/bizshuk/cc-plugin/pkg/topology"
	"github.com/spf13/cobra"
)

var (
	queryInboundOnly  bool
	queryOutboundOnly bool
	queryDepth        int
)

// QueryCmd lists inbound and outbound edges for an entity.
var QueryCmd = &cobra.Command{
	Use:   "query <entity>",
	Short: "List inbound and outbound edges for an entity",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if queryInboundOnly && queryOutboundOnly {
			return fmt.Errorf("--in and --out cannot be used together")
		}
		if queryDepth < 1 || queryDepth > 2 {
			return fmt.Errorf("depth must be 1 or 2")
		}
		topo, err := loadFromFlags(cmd)
		if err != nil {
			return err
		}
		if _, ok := topo.Entities[args[0]]; !ok {
			return fmt.Errorf("unknown entity %q", args[0])
		}
		for _, line := range queryEdges(topo, args[0], queryDepth, queryInboundOnly, queryOutboundOnly) {
			fmt.Fprintln(cmd.OutOrStdout(), line)
		}
		return nil
	},
}

func init() {
	QueryCmd.Flags().BoolVar(&queryInboundOnly, "in", false, "show inbound edges only")
	QueryCmd.Flags().BoolVar(&queryOutboundOnly, "out", false, "show outbound edges only")
	QueryCmd.Flags().IntVar(&queryDepth, "depth", 1, "query depth (1 or 2)")
}

func queryEdges(topo *topologypkg.Topology, entity string, depth int, inboundOnly, outboundOnly bool) []string {
	frontier := map[string]bool{entity: true}
	seenEntities := map[string]bool{entity: true}
	seenEdges := map[string]bool{}
	var lines []string

	for level := 0; level < depth; level++ {
		next := map[string]bool{}
		for _, edge := range topo.Edges() {
			includeOutbound := !inboundOnly && frontier[edge.FromEntity]
			includeInbound := !outboundOnly && frontier[edge.ToEntity]
			if !includeOutbound && !includeInbound {
				continue
			}
			line := formatEdge(edge)
			if !seenEdges[line] {
				seenEdges[line] = true
				lines = append(lines, line)
			}
			for _, candidate := range []string{edge.FromEntity, edge.ToEntity} {
				if !seenEntities[candidate] {
					next[candidate] = true
					seenEntities[candidate] = true
				}
			}
		}
		frontier = next
	}
	sort.Strings(lines)
	return lines
}

func formatEdge(edge topologypkg.TopoEdge) string {
	from := edge.FromEntity
	if edge.FromDim != "" {
		from += "#" + edge.FromDim
	}
	to := edge.ToEntity
	if edge.ToDim != "" {
		to += "#" + edge.ToDim
	}
	return fmt.Sprintf("%s -%s-> %s", from, edge.Relation, to)
}
