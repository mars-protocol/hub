package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	icatypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/types"

	"github.com/mars-protocol/hub/x/shuttle/types"
)

type queryServer struct{ k Keeper }

// NewQuerySErverImpl creates an implementation of the `QueryServer` interface
// for the given keeper.
func NewQueryServerImpl(k Keeper) types.QueryServer {
	return &queryServer{k}
}

func (qs queryServer) Account(goCtx context.Context, req *types.QueryAccountRequest) (*types.QueryAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner := qs.k.GetModuleAddress()
	portID, err := icatypes.NewControllerPortID(owner.String())
	if err != nil {
		return nil, err
	}

	address, found := qs.k.icaControllerKeeper.GetInterchainAccountAddress(ctx, req.ConnectionId, portID)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrAccountNotFound, "no interchain account exists on %s", req.ConnectionId)
	}

	return &types.QueryAccountResponse{Address: address}, nil
}

func (qs queryServer) Accounts(goCtx context.Context, req *types.QueryAccountsRequest) (*types.QueryAccountsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner := qs.k.GetModuleAddress()
	portID, err := icatypes.NewControllerPortID(owner.String())
	if err != nil {
		return nil, err
	}

	// the icaControllerKeeper does not provide a method to enumerate interchain
	// accounts of a particular owner or port id. instead, we simply fetch _all_
	// interchain accounts and filter them by the port id.
	allAccounts := qs.k.icaControllerKeeper.GetAllInterchainAccounts(ctx)
	accounts := []types.QueryAccountsResponseItem{}
	for _, account := range allAccounts {
		if account.PortId == portID {
			accounts = append(accounts, types.QueryAccountsResponseItem{
				ConnectionId: account.ConnectionId,
				Address:      account.AccountAddress,
			})
		}
	}

	return &types.QueryAccountsResponse{Accounts: accounts}, nil
}
