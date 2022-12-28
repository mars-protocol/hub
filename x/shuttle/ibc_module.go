package shuttle

import (
	"encoding/hex"

	"github.com/gogo/protobuf/proto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	ibcchanneltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	ibcporttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v6/modules/core/exported"

	"github.com/CosmWasm/wasmd/x/wasm"

	"github.com/mars-protocol/hub/x/shuttle/keeper"
)

// IBCModule implements the ICS26 interface for the shuttle module.
type IBCModule struct{ k keeper.Keeper }

// NewIBCModule creates a new IBCModule given the keeper.
func NewIBCModule(k keeper.Keeper) ibcporttypes.IBCModule {
	return IBCModule{k}
}

func (im IBCModule) OnChanOpenInit(
	ctx sdk.Context,
	order ibcchanneltypes.Order,
	connectionHops []string,
	portID string,
	channelID string,
	chanCap *capabilitytypes.Capability,
	counterparty ibcchanneltypes.Counterparty,
	version string,
) (string, error) {
	// the module is supposed to validate the version string here.
	// however, the version string is provided by the module's msgServer as an
	// empty string. therefore we skip the validation here.
	return version, nil
}

func (im IBCModule) OnChanOpenTry(
	ctx sdk.Context,
	order ibcchanneltypes.Order,
	connectionHops []string,
	portID,
	channelID string,
	chanCap *capabilitytypes.Capability,
	counterparty ibcchanneltypes.Counterparty,
	counterpartyVersion string,
) (string, error) {
	// ICA channel handshake cannot be initiated from the host chain.
	// the controller middleware should have rejected this request.
	// if not, something seriously wrong must have happened, e.g. modules wired
	// incorrectly in app.go. we panic in this case.
	panic("UNREACHABLE: shuttle module OnChanOpenTry")
}

func (im IBCModule) OnChanOpenAck(
	ctx sdk.Context,
	portID,
	channelID string,
	counterpartyChannelID string,
	counterpartyVersion string,
) error {
	// counterpartyVersion is already validated by the controller middleware.
	// we assume it's valid and don't validate again here.
	return nil
}

func (im IBCModule) OnChanOpenConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	// see the comment in OnChanOpenTry on why we panic here
	panic("UNREACHABLE: shuttle module OnChanOpenConfirm")
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
	// the ICA controller does not expect to receive any packet.
	// the controller middleware should have rejected this request.
	panic("UNREACHABLE: shuttle module OnRecvPacket")
}

// OnAcknowledgementPacket parses the acknowledgement and prints log messages.
//
// Although the shuttle module can send both ICS-20 and ICS-27 packets, only
// ICS-27 acknowledgements are routed here. ICS-20 packets are handled by the
// ibctransfer module alone.
//
// This function is mostly copied from interchain-account-demo:
// https://github.com/cosmos/interchain-accounts-demo/blob/v0.4.3/x/inter-tx/ibc_module.go#L108
func (im IBCModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet ibcchanneltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	logger := im.k.Logger(ctx)
	logger.Info(
		"received ICS-27 packet acknowledgement",
		"channel", packet.DestinationChannel,
		"sequence", packet.Sequence,
	)

	var ack ibcchanneltypes.Acknowledgement
	if err := ibcchanneltypes.SubModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-27 packet acknowledgement: %v", err)
	}

	var txMsgData sdk.TxMsgData
	if err := proto.Unmarshal(ack.GetResult(), &txMsgData); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-27 tx message data: %v", err)
	}

	switch len(txMsgData.Data) {
	// sdk 0.46
	case 0:
		for _, msgResp := range txMsgData.GetMsgResponses() {
			logger.Info(
				"message response in ICS-27 packet",
				"msgType", msgResp.TypeUrl,
				"response", msgResp.GoString(),
			)
		}
		return nil

	// sdk 0.45 or below
	default:
		for _, msgData := range txMsgData.Data {
			response, err := handleMsgData(msgData)
			if err != nil {
				return err
			}

			logger.Info(
				"message response in ICS-27 packet response",
				"msgType", msgData.MsgType,
				"response", response,
			)
		}
		return nil
	}
}

// handleMsgData parses the message response and return a human readable string
// representing the response data, for use in logging.
//
// We can't support every existing message types out there. Instead we only
// support a selected few:
//
//   - bank: MsgSend
//   - staking: MsgDelegate, MsgUndelegate
//   - distribution: MsgWithdrawRewards
//   - wasm: MsgExecuteContract, MsgMigrateContract
//
// For the others, we simply return the hex encoded bytes.
func handleMsgData(msgData *sdk.MsgData) (string, error) { //nolint:staticcheck // This function parses data from sdk 0.45 chains, so of course it contains deprecated stuff. Not my problem lol
	switch msgData.MsgType {
	case sdk.MsgTypeURL(&banktypes.MsgSend{}):
		return handleProtoMsg[*banktypes.MsgSendResponse](msgData.Data, "bank/MsgSend")
	case sdk.MsgTypeURL(&stakingtypes.MsgDelegate{}):
		return handleProtoMsg[*stakingtypes.MsgDelegateResponse](msgData.Data, "staking/MsgDelegate")
	case sdk.MsgTypeURL(&stakingtypes.MsgUndelegate{}):
		return handleProtoMsg[*stakingtypes.MsgUndelegateResponse](msgData.Data, "staking/MsgUndelegate")
	case sdk.MsgTypeURL(&distrtypes.MsgWithdrawDelegatorReward{}):
		return handleProtoMsg[*distrtypes.MsgWithdrawDelegatorRewardResponse](msgData.Data, "distr/MsgWithdrawDelegatorReward")
	case sdk.MsgTypeURL(&wasm.MsgExecuteContract{}):
		return handleProtoMsg[*wasm.MsgExecuteContractResponse](msgData.Data, "wasm/MsgExecuteContract")
	case sdk.MsgTypeURL(&wasm.MsgMigrateContract{}):
		return handleProtoMsg[*wasm.MsgMigrateContractResponse](msgData.Data, "wasm/MsgMigrateContract")
	default:
		return hex.EncodeToString(msgData.Data), nil
	}
}

// handleProtoMsg unmarshals bytes into a given proto message and stringifies it.
func handleProtoMsg[T proto.Message](bz []byte, tyName string) (string, error) {
	var msg T
	if err := proto.Unmarshal(bz, msg); err != nil {
		return "", sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "cannot unmarshal %s response: %v", tyName, err)
	}

	return msg.String(), nil
}

// OnTimeoutPacket prints a log message indicating the packet has timed out.
func (im IBCModule) OnTimeoutPacket(
	ctx sdk.Context,
	packet ibcchanneltypes.Packet,
	relayer sdk.AccAddress,
) error {
	im.k.Logger(ctx).Info(
		"ICS-27 packet timed out",
		"channel", packet.DestinationChannel,
		"sequence", packet.Sequence,
	)

	return nil
}
