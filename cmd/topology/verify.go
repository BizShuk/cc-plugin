package topology

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

// VerifyCmd returns the graph-integrity verification command.
func VerifyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "verify",
		Short: "Check links, kinds, backlinks, and duplicate entities",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			topo, err := loadFromFlags(cmd)
			if err != nil {
				return err
			}
			findings := topo.Verify()
			if len(findings) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "OK")
				return nil
			}
			for _, finding := range findings {
				fmt.Fprintln(cmd.OutOrStdout(), finding)
			}
			return errors.New("topology verification failed")
		},
	}
}
