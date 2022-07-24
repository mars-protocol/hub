package client

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govclientcli "github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govclientrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/mars-protocol/hub/x/safetyfund/types"
)

// ProposalHandler is the safety fund spend proposal handler
var ProposalHandler = govclient.NewProposalHandler(getCmdSubmitProposal, getProposalRESTHandler)

func getCmdSubmitProposal() *cobra.Command {
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

			title, err := cmd.Flags().GetString(govclientcli.FlagTitle)
			if err != nil {
				return fmt.Errorf("invalid title: %s", err)
			}

			description, err := cmd.Flags().GetString(govclientcli.FlagDescription)
			if err != nil {
				return fmt.Errorf("invalid description: %s", err)
			}

			depositStr, err := cmd.Flags().GetString(govclientcli.FlagDeposit)
			if err != nil {
				return fmt.Errorf("invalid deposit: %s", err)
			}

			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return fmt.Errorf("invalid deposit: %s", err)
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

	cmd.Flags().String(govclientcli.FlagTitle, "", "Title of proposal")
	cmd.Flags().String(govclientcli.FlagDescription, "", "Description of proposal")
	cmd.Flags().String(govclientcli.FlagDeposit, "", "Deposit of proposal")

	return cmd
}

func getProposalRESTHandler(clientCtx client.Context) govclientrest.ProposalRESTHandler {
	return govclientrest.ProposalRESTHandler{
		SubRoute: "safety_fund_spend",
		Handler:  func(w http.ResponseWriter, r *http.Request) {}, // deprecated, do nothing
	}
}
