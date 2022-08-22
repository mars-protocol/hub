package shuttle

import (
	"encoding/hex"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/gogo/protobuf/proto"

	ibcchanneltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	ibcporttypes "github.com/cosmos/ibc-go/v3/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	ibcexported "github.com/cosmos/ibc-go/v3/modules/core/exported"

	wasm "github.com/CosmWasm/wasmd/x/wasm"

	"github.com/mars-protocol/hub/x/shuttle/keeper"
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
) error {
	// claim channel capability passed back by IBC module
	return im.ClaimCapability(ctx, channelCap, ibchost.ChannelCapabilityPath(portID, channelID))
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
	return "", nil
}

func (im IBCModule) OnChanOpenAck(
	ctx sdk.Context,
	portID,
	channelID string,
	counterpartyChannelID string,
	counterpartyVersion string,
) error {
	return nil
}

func (im IBCModule) OnChanOpenConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return nil
}

func (im IBCModule) OnChanCloseInit(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return nil
}

func (im IBCModule) OnChanCloseConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return nil
}

func (im IBCModule) OnRecvPacket(
	ctx sdk.Context,
	packet ibcchanneltypes.Packet,
	relayer sdk.AccAddress,
) ibcexported.Acknowledgement {
	return ibcchanneltypes.NewErrorAcknowledgement("cannot receive packet via interchain accounts authentication module")
}

func (im IBCModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet ibcchanneltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	var ack ibcchanneltypes.Acknowledgement
	if err := ibcchanneltypes.SubModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "failed to unmarshal ICS-27 packet acknowledgement: %v", err)
	}

	txMsgData := &sdk.TxMsgData{}
	if err := proto.Unmarshal(ack.GetResult(), txMsgData); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "failed to unmarshal ICS-27 tx message data: %v", err)
	}

	switch len(txMsgData.Data) {
	case 0:
		// TODO: handle for sdk 0.46.x
		return nil
	default:
		for _, msgData := range txMsgData.Data {
			response, err := handleMsgData(ctx, msgData)
			if err != nil {
				return err
			}

			im.Logger(ctx).Info("message response in ICS-27 packet response", "response", response)
		}
		return nil
	}
}

func handleMsgData(ctx sdk.Context, msgData *sdk.MsgData) (string, error) {
	switch msgData.MsgType {
	case sdk.MsgTypeURL(&wasm.MsgExecuteContract{}):
		msgResponse := &wasm.MsgExecuteContractResponse{}
		if err := proto.Unmarshal(msgData.Data, msgResponse); err != nil {
			return "", sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "failed to unmarshal response for MsgExecuteContract: %s", err.Error())
		}

		return msgResponse.String(), nil

	case sdk.MsgTypeURL(&wasm.MsgMigrateContract{}):
		msgResponse := &wasm.MsgMigrateContractResponse{}
		if err := proto.Unmarshal(msgData.Data, msgResponse); err != nil {
			return "", sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "failed to unmarshal response for MsgMigrateContract: %s", err.Error())
		}

		return msgResponse.String(), nil

	default:
		return hex.EncodeToString(msgData.Data), nil
	}
}

func (im IBCModule) OnTimeoutPacket(
	ctx sdk.Context,
	packet ibcchanneltypes.Packet,
	relayer sdk.AccAddress,
) error {
	// TODO: if the packet times out, we should save the msg somewhere and let anyone execute it
	// later once the congestion clears up, instead of having to do the same gov poll again
	return nil
}
