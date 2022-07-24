package client

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govclientcli "github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	marsutils "github.com/mars-protocol/hub/utils"

	"github.com/mars-protocol/hub/x/incentives/types"
)

var (
	CreateIncentivesProposalHandler    = govclient.NewProposalHandler(getCreateIncentivesProposalCmd, marsutils.GetProposalRESTHandler("create_incentives_schedule"))
	TerminateIncentivesProposalHandler = govclient.NewProposalHandler(getTerminateIncentivesProposalCmd, marsutils.GetProposalRESTHandler("terminate_incentives_schedule"))
)

func getCreateIncentivesProposalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-incentives-schedule [start-time] [end-time] [amount] --title [text] --description [text] --deposit [amount]",
		Args:  cobra.ExactArgs(3),
		Short: "Submit a create incentives schedule proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			startTime, err := time.Parse(time.RFC3339, args[0])
			if err != nil {
				return fmt.Errorf("invalid start time: %s", err)
			}

			endTime, err := time.Parse(time.RFC3339, args[1])
			if err != nil {
				return fmt.Errorf("invalid end time: %s", err)
			}

			amount, err := sdk.ParseCoinsNormalized(args[2])
			if err != nil {
				return fmt.Errorf("invalid amount: %s", err)
			}

			title, description, deposit, err := marsutils.ParseGovProposalFlags(cmd)
			if err != nil {
				return err
			}

			proposal := &types.CreateIncentivesScheduleProposal{
				Title:       title,
				Description: description,
				StartTime:   startTime,
				EndTime:     endTime,
				Amount:      amount,
			}

			if err := proposal.ValidateBasic(); err != nil {
				return fmt.Errorf("invalid proposal: %s", err)
			}

			msg, err := govtypes.NewMsgSubmitProposal(proposal, deposit, clientCtx.GetFromAddress())
			if err != nil {
				return fmt.Errorf("failed to create msg: %s", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(govclientcli.FlagTitle, "", "Title of proposal")
	cmd.Flags().String(govclientcli.FlagDescription, "", "Description of proposal")
	cmd.Flags().String(govclientcli.FlagDeposit, "", "Deposit of proposal")

	return cmd
}

func getTerminateIncentivesProposalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "terminate-incentives-schedules [ids] --title [text] --description [text] --deposit [amount]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a terminate incentives proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			idStrs := strings.Split(args[0], ",")
			ids := []uint64{}
			for _, idStr := range idStrs {
				id, err := strconv.ParseUint(idStr, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid ids: %s", err)
				}

				ids = append(ids, id)
			}

			title, description, deposit, err := marsutils.ParseGovProposalFlags(cmd)
			if err != nil {
				return err
			}

			proposal := &types.TerminateIncentivesSchedulesProposal{
				Title:       title,
				Description: description,
				Ids:         ids,
			}

			if err := proposal.ValidateBasic(); err != nil {
				return fmt.Errorf("invalid proposal: %s", err)
			}

			msg, err := govtypes.NewMsgSubmitProposal(proposal, deposit, clientCtx.GetFromAddress())
			if err != nil {
				return fmt.Errorf("failed to create msg: %s", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(govclientcli.FlagTitle, "", "Title of proposal")
	cmd.Flags().String(govclientcli.FlagDescription, "", "Description of proposal")
	cmd.Flags().String(govclientcli.FlagDeposit, "", "Deposit of proposal")

	return cmd
}
