package memory

import (
	"fmt"

	"github.com/spf13/cobra"
)

func ResetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reset",
		Short: "Reset distilled data status (clear sqlite state tables: cursor, seen, and distilled)",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := NewStateStore()
			if err != nil {
				return err
			}
			defer store.Close()

			if err := store.Reset(); err != nil {
				return err
			}

			fmt.Println("Distilled data status has been successfully reset.")
			return nil
		},
	}

	return cmd
}
