package shuttle

import (
	"encoding/hex"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/gogo/protobuf/proto"

	ibcchanneltypes "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"
	ibcporttypes "github.com/cosmos/ibc-go/v4/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/v4/modules/core/24-host"
	ibcexported "github.com/cosmos/ibc-go/v4/modules/core/exported"

	wasm "github.com/CosmWasm/wasmd/x/wasm"

	"github.com/mars-protocol/hub/x/shuttle/keeper"
	"github.com/mars-protocol/hub/x/shuttle/types"
)

var _ ibcporttypes.IBCModule = IBCModule{}

// IBCModule implements the ICS26 interface for interchain accounts controller chains
type IBCModule struct{ keeper.Keeper }

// NewIBCModule creates a new IBCModule given the keeper
func NewIBCModule(k keeper.Keeper) IBCModule {
	return IBCModule{k}
}

func (im IBCModule) OnChanOpenInit(
	ctx sdk.Context,
	order ibcchanneltypes.Order,
	connectionHops []string,
	portID string,
	channelID string,
	channelCap *capabilitytypes.Capability,
	counterparty ibcchanneltypes.Counterparty,
	version string,
) (string, error) {
	im.Keeper.Logger(ctx).Info(
		"OnChanOpenInit",
		"portID", portID,
		"channelID", channelID,
		"counterpartyPortID", counterparty.PortId,
		"counterpartyChannelID", counterparty.ChannelId,
		"version", version,
	)

	// claim channel capability passed back by IBC module
	if err := im.ClaimCapability(ctx, channelCap, ibchost.ChannelCapabilityPath(portID, channelID)); err != nil {
		return "", err
	}

	im.Keeper.Logger(ctx).Info(
		"Claimed channel capability",
		"module", types.ModuleName,
		"capability", channelCap.String(),
	)

	// Since ibc-go v4, the `OnChanOpenInit` function needs to validate the
	// version string, and if it's valid, include it in the return values.
	//
	// From the comments of the `IBCModule` interface:
	//
	// - If the provided version string is non-empty, OnChanOpenInit should
	//   return the version string if valid or an error if the provided version
	//   is invalid.
	// - If the version string is empty, OnChanOpenInit is expected to
	//   return a default version string representing the version(s) it supports.
	// - If there is no default version string for the application,
	//   it should return an error if provided version is empty string.
	//
	// Here we attempt to parse the version string into `icatypes.Metadata` and
	// do some basic validations.
	//
	// TODO: validate version

	return version, nil
}

func (im IBCModule) OnChanOpenTry(
	ctx sdk.Context,
	order ibcchanneltypes.Order,
	connectionHops []string,
	portID,
	channelID string,
	channelCap *capabilitytypes.Capability,
	counterparty ibcchanneltypes.Counterparty,
	counterpartyVersion string,
) (version string, err error) {
	// we don't expect the host chain to open a channel with the shuttle module
	return "", types.ErrUnexpectedChannelOpen
}

func (im IBCModule) OnChanOpenAck(
	ctx sdk.Context,
	portID,
	channelID string,
	counterpartyChannelID string,
	counterpartyVersion string,
) error {
	im.Keeper.Logger(ctx).Info(
		"OnChanOpenInit",
		"portID", portID,
		"channelID", channelID,
		"counterpartyChannelID", counterpartyChannelID,
		"counterpartyVersion", counterpartyVersion,
	)

	// TODO: validate counterparty version

	return nil
}

func (im IBCModule) OnChanOpenConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	// we don't expect the host chain to open a channel with the shuttle module
	return types.ErrUnexpectedChannelOpen
}

func (im IBCModule) OnChanCloseInit(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	// we don't want to close the channel under any circumstance
	return types.ErrUnexpectedChannelClose
}

func (im IBCModule) OnChanCloseConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	// we don't want to close the channel under any circumstance
	return types.ErrUnexpectedChannelClose
}

func (im IBCModule) OnRecvPacket(
	ctx sdk.Context,
	packet ibcchanneltypes.Packet,
	relayer sdk.AccAddress,
) ibcexported.Acknowledgement {
	return ibcchanneltypes.NewErrorAcknowledgement(sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "interchain account auth module does not expect to receive any packet"))
}

func (im IBCModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet ibcchanneltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	im.Keeper.Logger(ctx).Info("OnAcknowledgementPacket", "packet", packet.String(), "relayer", relayer)

	var ack ibcchanneltypes.Acknowledgement
	if err := ibcchanneltypes.SubModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "failed to unmarshal ICS-27 packet acknowledgement: %v", err)
	}

	im.Keeper.Logger(ctx).Info("Successfully unmarshalled acknowledgement", "ack", ack.String())

	txMsgData := &sdk.TxMsgData{}
	if err := proto.Unmarshal(ack.GetResult(), txMsgData); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "failed to unmarshal ICS-27 tx message data: %v", err)
	}

	for _, msgData := range txMsgData.Data {
		msgType, msgResponse, err := handleMsgData(ctx, msgData)
		if err != nil {
			return err
		}

		im.Logger(ctx).Info("Message response in ICS-27 packet response", "msgType", msgType, "msgResponse", msgResponse)
	}
	return nil

}

// handleMsgData parses the msg response data for logging purpose.
func handleMsgData(ctx sdk.Context, msgData *sdk.MsgData) (msgType, msgResponse string, err error) {
	switch msgData.MsgType {
	case sdk.MsgTypeURL(&wasm.MsgExecuteContract{}):
		msgResponse := &wasm.MsgExecuteContractResponse{}
		if err := proto.Unmarshal(msgData.Data, msgResponse); err != nil {
			return "", "", sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "failed to unmarshal response for MsgExecuteContract: %s", err.Error())
		}

		return "MsgExecuteContract", msgResponse.String(), nil

	case sdk.MsgTypeURL(&wasm.MsgMigrateContract{}):
		msgResponse := &wasm.MsgMigrateContractResponse{}
		if err := proto.Unmarshal(msgData.Data, msgResponse); err != nil {
			return "", "", sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "failed to unmarshal response for MsgMigrateContract: %s", err.Error())
		}

		return "MsgMigrateContract", msgResponse.String(), nil

	default:
		return "unknown", hex.EncodeToString(msgData.Data), nil
	}
}

func (im IBCModule) OnTimeoutPacket(
	ctx sdk.Context,
	packet ibcchanneltypes.Packet,
	relayer sdk.AccAddress,
) error {
	im.Keeper.Logger(ctx).Info("OnTimeoutPacket", "packet", packet.String(), "relayer", relayer)

	return nil
}
