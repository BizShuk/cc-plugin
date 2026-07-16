package topology

import (
	"fmt"
	"os"
	"path/filepath"

	topologypkg "github.com/bizshuk/cc-plugin/pkg/topology"
	"github.com/spf13/cobra"
)

// BacklinksCmd returns the backlink regeneration command.
func BacklinksCmd() *cobra.Command {
	var write bool
	cmd := &cobra.Command{
		Use:   "backlinks",
		Short: "Recompute Backlinks sections from forward edges",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			topo, err := loadFromFlags(cmd)
			if err != nil {
				return err
			}
			return rewriteBacklinks(cmd, topo, write)
		},
	}
	cmd.Flags().BoolVar(&write, "write", false, "write changes to entity files")
	return cmd
}

// IndexCmd returns the topology-index regeneration command.
func IndexCmd() *cobra.Command {
	var write bool
	cmd := &cobra.Command{
		Use:   "index",
		Short: "Regenerate _index.md while preserving Frontier",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			topo, err := loadFromFlags(cmd)
			if err != nil {
				return err
			}
			return rewriteIndex(cmd, topo, write)
		},
	}
	cmd.Flags().BoolVar(&write, "write", false, "write changes to _index.md")
	return cmd
}

// RewriteCmd returns the combined backlink and index rewrite command.
func RewriteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rewrite",
		Short: "Rewrite backlinks and _index.md",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			topo, err := loadFromFlags(cmd)
			if err != nil {
				return err
			}
			if err := rewriteBacklinks(cmd, topo, true); err != nil {
				return err
			}
			return rewriteIndex(cmd, topo, true)
		},
	}
}

func rewriteBacklinks(cmd *cobra.Command, topo *topologypkg.Topology, write bool) error {
	updated := 0
	for _, name := range topo.Names() {
		entity := topo.Entities[name]
		raw, err := os.ReadFile(entity.Path)
		if err != nil {
			return fmt.Errorf("read %s: %w", entity.Path, err)
		}
		rendered := topologypkg.RenderBacklinksSection(string(raw), topo.BacklinksFor(name))
		if rendered == string(raw) {
			continue
		}
		updated++
		if write {
			if err := os.WriteFile(entity.Path, []byte(rendered), 0o644); err != nil {
				return fmt.Errorf("write %s: %w", entity.Path, err)
			}
		}
	}
	fmt.Fprintf(cmd.OutOrStdout(), "backlinks: %d file(s) %s\n", updated, writeAction(write))
	return nil
}

func rewriteIndex(cmd *cobra.Command, topo *topologypkg.Topology, write bool) error {
	path := filepath.Join(topo.Root, "_index.md")
	existing, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("read %s: %w", path, err)
	}
	rendered := topo.RenderIndex(string(existing))
	changed := rendered != string(existing)
	if changed && write {
		if err := os.WriteFile(path, []byte(rendered), 0o644); err != nil {
			return fmt.Errorf("write %s: %w", path, err)
		}
	}
	fmt.Fprintf(cmd.OutOrStdout(), "index: %t %s\n", changed, writeAction(write))
	return nil
}

func writeAction(write bool) string {
	if write {
		return "written"
	}
	return "would change"
}
