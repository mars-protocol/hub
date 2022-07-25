package client

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	marsutils "github.com/mars-protocol/hub/utils"

	"github.com/mars-protocol/hub/x/safetyfund/types"
)

// SafetyFundSpendProposalHandler is the safety fund spend proposal handler
var SafetyFundSpendProposalHandler = govclient.NewProposalHandler(getSafetyFundCommandProposalCmd, marsutils.GetProposalRESTHandler("safety_fund_spend"))

func getSafetyFundCommandProposalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "safety-fund-spend [recipient] [amount] --title [text] --description [text] --deposit [amount]",
		Args:  cobra.ExactArgs(2),
		Short: "Submit a safety fund spend proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			recipient := args[0]
			if _, err := sdk.AccAddressFromBech32(recipient); err != nil {
				return fmt.Errorf("invalid recipient: %s", err)
			}

			amount, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return fmt.Errorf("invalid amount: %s", err)
			}

			title, description, deposit, err := marsutils.ParseGovProposalFlags(cmd)
			if err != nil {
				return nil
			}

			proposal := &types.SafetyFundSpendProposal{
				Title:       title,
				Description: description,
				Recipient:   recipient,
				Amount:      amount,
			}

			if err := proposal.ValidateBasic(); err != nil {
				return fmt.Errorf("invalid proposal: %s", err)
			}

			msg, err := govtypes.NewMsgSubmitProposal(proposal, deposit, clientCtx.GetFromAddress())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	marsutils.AddGovProposalFlags(cmd)

	return cmd
}
