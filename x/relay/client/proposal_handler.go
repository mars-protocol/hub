package client

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"

	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	marsutils "github.com/mars-protocol/hub/utils"

	"github.com/mars-protocol/hub/x/relay/types"
)

var (
	ExecuteRemoteContractProposalHandler = govclient.NewProposalHandler(getExecuteRemoteContractProposalCmd, marsutils.GetProposalRESTHandler("execute_remote_contract"))
	MigrateRemoteContractProposalHandler = govclient.NewProposalHandler(getMigrateRemoteContractProposalCmd, marsutils.GetProposalRESTHandler("migrate_remote_contract"))
)

func getExecuteRemoteContractProposalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execute-remote-contract [connection-id] [contract-addr] [json-encoded-msg] --title [text] --description [text] --deposit [amount]",
		Args:  cobra.ExactArgs(3),
		Short: "Submit an execute remote wasm contract proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			title, description, deposit, err := marsutils.ParseGovProposalFlags(cmd)
			if err != nil {
				return err
			}

			proposal := &types.ExecuteRemoteContractProposal{
				Title:        title,
				Description:  description,
				ConnectionId: args[0],
				Contract:     args[1],
				ExecuteMsg:   []byte(args[2]),
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

	marsutils.AddGovProposalFlags(cmd)

	return cmd
}

func getMigrateRemoteContractProposalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate-remote-contract [connection-id] [contract-addr] [code-id] [json-encoded-msg] --title [text] --description [text] --deposit [amount]",
		Args:  cobra.ExactArgs(4),
		Short: "Submit a migrate remote wasm contract proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			title, description, deposit, err := marsutils.ParseGovProposalFlags(cmd)
			if err != nil {
				return err
			}

			codeId, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			proposal := &types.MigrateRemoteContractProposal{
				Title:        title,
				Description:  description,
				ConnectionId: args[0],
				Contract:     args[1],
				CodeId:       codeId,
				MigrateMsg:   []byte(args[3]),
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

	marsutils.AddGovProposalFlags(cmd)

	return cmd
}
