package utils

import (
	"fmt"

	"github.com/spf13/cobra"

	sdk "github.com/cosmos/cosmos-sdk/types"

	govclientcli "github.com/cosmos/cosmos-sdk/x/gov/client/cli"
)

// ParseGovProposalFlags parses flags related to creating governance proposals added to the given command
func ParseGovProposalFlags(cmd *cobra.Command) (title, description string, deposit sdk.Coins, err error) {
	title, err = cmd.Flags().GetString(govclientcli.FlagTitle)
	if err != nil {
		return "", "", sdk.NewCoins(), fmt.Errorf("invalid title: %s", err)
	}

	description, err = cmd.Flags().GetString(govclientcli.FlagDescription)
	if err != nil {
		return "", "", sdk.NewCoins(), fmt.Errorf("invalid description: %s", err)
	}

	depositStr, err := cmd.Flags().GetString(govclientcli.FlagDeposit)
	if err != nil {
		return "", "", sdk.NewCoins(), fmt.Errorf("invalid deposit: %s", err)
	}

	deposit, err = sdk.ParseCoinsNormalized(depositStr)
	if err != nil {
		return "", "", sdk.NewCoins(), fmt.Errorf("invalid deposit: %s", err)
	}

	return title, description, deposit, nil
}

// AddGovProposalFlags adds flags related to creating governance proposal to the given command
func AddGovProposalFlags(cmd *cobra.Command) {
	cmd.Flags().String(govclientcli.FlagTitle, "", "Title of proposal")
	cmd.Flags().String(govclientcli.FlagDescription, "", "Description of proposal")
	cmd.Flags().String(govclientcli.FlagDeposit, "", "Deposit of proposal")
}
