package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/mars-protocol/hub/x/safety/types"
)

// GetQueryCmd returns the parent command for all safety module query commands
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "safety-fund",
		Short:                      "Querying commands for the safety fund module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		getBalancesCmd(),
	)

	return cmd
}

func getBalancesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "balances",
		Short: "Query the amount of coins available in the safety fund",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Balances(cmd.Context(), &types.QueryBalancesRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
