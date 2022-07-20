package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
)

// GetTxCmd returns the parent command for all safetyfund module tx commands
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "safety-fund",
		Short:                      "Safety fund transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
	//...
	)

	return cmd
}
