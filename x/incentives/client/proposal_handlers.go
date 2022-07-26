package client

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	marsutils "github.com/mars-protocol/hub/utils"

	"github.com/mars-protocol/hub/x/incentives/types"
)

var (
	CreateIncentivesProposalHandler    = govclient.NewProposalHandler(getCreateIncentivesProposalCmd)
	TerminateIncentivesProposalHandler = govclient.NewProposalHandler(getTerminateIncentivesProposalCmd)
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

			msg, err := govv1.NewMsgSubmitProposal(proposal, deposit, clientCtx.GetFromAddress())
			if err != nil {
				return fmt.Errorf("failed to create msg: %s", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	marsutils.AddGovProposalFlags(cmd)

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

			ids, err := marsutils.StringToUintArray(args[0], ",")
			if err != nil {
				return err
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

			msg, err := govv1.NewMsgSubmitProposal(proposal, deposit, clientCtx.GetFromAddress())
			if err != nil {
				return fmt.Errorf("failed to create msg: %s", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	marsutils.AddGovProposalFlags(cmd)

	return cmd
}
