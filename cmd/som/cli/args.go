package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ExactArgs returns an error if there are not exactly n args.
func ExactArgs(n int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != n {
			return fmt.Errorf("command '%s' accepts %d arg(s), received %d", cmd.Name(), n, len(args))
		}
		return nil
	}
}
