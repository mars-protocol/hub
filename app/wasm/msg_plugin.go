package wasm

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	customdistrkeeper "github.com/mars-protocol/hub/x/distribution/keeper"
)

func CustomMessageDecorator(distrKeeper *customdistrkeeper.Keeper) func(wasmkeeper.Messenger) wasmkeeper.Messenger {
	return func(old wasmkeeper.Messenger) wasmkeeper.Messenger {
		return &CustomMessenger{
			wrapped:     old,
			distrKeeper: distrKeeper,
		}
	}
}

type CustomMessenger struct {
	wrapped     wasmkeeper.Messenger
	distrKeeper *customdistrkeeper.Keeper
}

// CustomKeeper must implement the `wasmkeeper.Messenger` interface
var _ wasmkeeper.Messenger = (*CustomMessenger)(nil)

func (m *CustomMessenger) DispatchMsg(
	ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg,
) ([]sdk.Event, [][]byte, error) {
	// if the msg is a custom msg, parse it into `MarsMsg` then dispatch to the appropriate Mars module
	// otherwise, simply dispatch it to the wrapped messenger
	if msg.Custom != nil {
		var marsMsg MarsMsg
		if err := json.Unmarshal(msg.Custom, &marsMsg); err != nil {
			return nil, nil, sdkerrors.Wrapf(err, "invalid custom msg: %s", msg.Custom)
		}

		if marsMsg.FundCommunityPool != nil {
			return fundCommunityPool(ctx, m.distrKeeper, contractAddr, marsMsg.FundCommunityPool)
		}
	}

	return m.wrapped.DispatchMsg(ctx, contractAddr, contractIBCPortID, msg)
}

func fundCommunityPool(
	ctx sdk.Context, k *customdistrkeeper.Keeper, contractAddr sdk.AccAddress,
	fundCommunityPool *FundCommunityPool,
) ([]sdk.Event, [][]byte, error) {
	msgServer := distrkeeper.NewMsgServerImpl(k.Keeper)

	msg := &distrtypes.MsgFundCommunityPool{
		Amount:    fundCommunityPool.Amount,
		Depositor: contractAddr.String(),
	}

	if _, err := msgServer.FundCommunityPool(sdk.WrapSDKContext(ctx), msg); err != nil {
		return nil, nil, err
	}

	return nil, nil, nil
}
