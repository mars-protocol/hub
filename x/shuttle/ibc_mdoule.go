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

var _ ibcporttypes.IBCModule = IBCModule{}

// IBCModule implements the ICS26 interface for the shuttle module.
type IBCModule struct{ k keeper.Keeper }

// NewIBCModule creates a new IBCModule given the keeper.
func NewIBCModule(k keeper.Keeper) IBCModule {
	return IBCModule{k}
}

// OnChanOpenInit implements the IBCModule interface
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
	return version, nil
}

// OnChanOpenTry implements the IBCModule interface
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
	return "", nil
}

// OnChanOpenAck implements the IBCModule interface
func (im IBCModule) OnChanOpenAck(
	ctx sdk.Context,
	portID,
	channelID string,
	counterpartyChannelID string,
	counterpartyVersion string,
) error {
	return nil
}

// OnChanOpenConfirm implements the IBCModule interface
func (im IBCModule) OnChanOpenConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return nil
}

// OnChanCloseInit implements the IBCModule interface
func (im IBCModule) OnChanCloseInit(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return nil
}

// OnChanCloseConfirm implements the IBCModule interface
func (im IBCModule) OnChanCloseConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return nil
}

// OnRecvPacket implements the IBCModule interface
func (im IBCModule) OnRecvPacket(
	ctx sdk.Context,
	packet ibcchanneltypes.Packet,
	relayer sdk.AccAddress,
) ibcexported.Acknowledgement {
	// TODO
	return ibcchanneltypes.NewErrorAcknowledgement(errors.New("UNIMPLEMENTED"))
}

// OnAcknowledgementPacket implements the IBCModule interface
func (im IBCModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet ibcchanneltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	// TODO
	return errors.New("UNIMPLEMENTED")
}

// OnTimeoutPacket implements the IBCModule interface.
func (im IBCModule) OnTimeoutPacket(
	ctx sdk.Context,
	packet ibcchanneltypes.Packet,
	relayer sdk.AccAddress,
) error {
	// TODO
	return errors.New("UNIMPLEMENTED")
}
