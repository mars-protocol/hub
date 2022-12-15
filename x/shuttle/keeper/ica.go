package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	icatypes "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"
	ibchost "github.com/cosmos/ibc-go/v4/modules/core/24-host"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	"github.com/mars-protocol/hub/x/shuttle/types"
)

func (k Keeper) RegisterAccount(ctx sdk.Context, connectionID string) error {
	macc := k.GetModuleAddress().String()
	portID, err := icatypes.NewControllerPortID(macc)
	if err != nil {
		return err
	}

	// query the counterparty connection id, which is used in the ICA metadata
	// (see below)
	connectionEnd, found := k.ibcKeeper.ConnectionKeeper.GetConnection(ctx, connectionID)
	if !found {
		return sdkerrors.Wrapf(types.ErrConnectionNotFound, "connection not found for id %s", connectionID)
	}

	// TODO: if a channel is being created (doing handshakes, not yet completed) then this check
	// will pass, but we don't want it to pass in this situation
	_, found = k.icaControllerKeeper.GetInterchainAccountAddress(ctx, connectionID, portID)
	if found {
		return sdkerrors.Wrapf(types.ErrAccountExists, "interchain account already exists for %s", connectionID)
	}

	// Since ibc-go v4, it is required that when registering an interchain
	// account, a JSON string which includes necessary metadata is included.
	//
	// See migration guide:
	// https://github.com/cosmos/ibc-go/blob/v4.2.0/docs/migrations/v3-to-v4.md
	icaMetadata := icatypes.Metadata{
		Version:                icatypes.Version,
		ControllerConnectionId: connectionID,
		HostConnectionId:       connectionEnd.Counterparty.ConnectionId,
		Encoding:               icatypes.EncodingProtobuf,
		TxType:                 icatypes.TxTypeSDKMultiMsg,
	}

	appVersion, err := icatypes.ModuleCdc.MarshalJSON(&icaMetadata)
	if err != nil {
		return err
	}

	return k.icaControllerKeeper.RegisterInterchainAccount(ctx, connectionID, macc, string(appVersion))
}

// ExecuteRemoteContract sends a message to execute a remote contract on the chain specified by the
// given connection ID via interchain account
func (k Keeper) ExecuteRemoteContract(ctx sdk.Context, connectionID, contract string, msg wasmtypes.RawContractMessage, funds sdk.Coins) (uint64, error) {
	return k.sendInterchainTx(
		ctx,
		connectionID,
		func(sender string) []sdk.Msg {
			return []sdk.Msg{&wasmtypes.MsgExecuteContract{
				Sender:   sender,
				Contract: contract,
				Msg:      msg,
				Funds:    funds,
			}}
		},
	)
}

// MigrateRemoteContract sends a message to migrate a remote contract on the chain specified by the
// given connection ID via interchain account
func (k Keeper) MigrateRemoteContract(ctx sdk.Context, connectionID, contract string, codeID uint64, msg wasmtypes.RawContractMessage) (uint64, error) {
	return k.sendInterchainTx(
		ctx,
		connectionID,
		func(sender string) []sdk.Msg {
			return []sdk.Msg{&wasmtypes.MsgMigrateContract{
				Sender:   sender,
				Contract: contract,
				CodeID:   codeID,
				Msg:      msg,
			}}
		},
	)
}

func (k Keeper) sendInterchainTx(ctx sdk.Context, connectionID string, buildMsgs func(sender string) []sdk.Msg) (uint64, error) {
	macc := k.GetModuleAddress().String()
	portID, err := icatypes.NewControllerPortID(macc)
	if err != nil {
		return 0, err
	}

	icacc, found := k.icaControllerKeeper.GetInterchainAccountAddress(ctx, connectionID, portID)
	if !found {
		return 0, sdkerrors.Wrapf(types.ErrAccountNotFound, "interchain account does not exist for %s", connectionID)
	}

	channelID, found := k.icaControllerKeeper.GetActiveChannelID(ctx, connectionID, portID)
	if !found {
		return 0, sdkerrors.Wrapf(icatypes.ErrActiveChannelNotFound, "failed to retrieve active channel for port %s", portID)
	}

	chanCap, found := k.scopedKeeper.GetCapability(ctx, ibchost.ChannelCapabilityPath(portID, channelID))
	if !found {
		return 0, sdkerrors.Wrap(ibcchanneltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	data, err := icatypes.SerializeCosmosTx(k.cdc, buildMsgs(icacc))
	if err != nil {
		return 0, err
	}

	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: data,
	}

	params := k.GetParams(ctx)
	timeoutTimestamp := time.Now().Add(params.TimeoutDuration).UnixNano()

	return k.icaControllerKeeper.SendTx(ctx, chanCap, connectionID, portID, packetData, uint64(timeoutTimestamp))
}
