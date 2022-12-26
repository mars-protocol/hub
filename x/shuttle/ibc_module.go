package shuttle

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"

	ibcchanneltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	ibcporttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v6/modules/core/exported"

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
	// we don't want to ever close the ICA channel, and provide no message type
	// for doing so.
	// if this function is somehow invoked, then something seriously wrong has
	// happened. we panic in this case.
	panic("UNREACHABLE: shuttle module OnChanCloseInit")
}

func (im IBCModule) OnChanCloseConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	// see the comment in OnChanCloseInit on why we panic here
	panic("UNREACHABLE: shuttle module OnChanCloseConfirm")
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

func (im IBCModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet ibcchanneltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	// TODO
	return errors.New("UNIMPLEMENTED")
}

func (im IBCModule) OnTimeoutPacket(
	ctx sdk.Context,
	packet ibcchanneltypes.Packet,
	relayer sdk.AccAddress,
) error {
	// TODO
	return errors.New("UNIMPLEMENTED")
}
