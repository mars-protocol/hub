package relay

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"

	ibcchanneltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	ibcporttypes "github.com/cosmos/ibc-go/v3/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v3/modules/core/exported"

	"github.com/mars-protocol/hub/x/relay/keeper"
)

var _ ibcporttypes.IBCModule = IBCModule{}

// IBCModule implements the ICS26 interface for interchain accounts controller chains
type IBCModule struct {
	keeper keeper.Keeper
}

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
	return errors.New("unimplemented")
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
	return "", errors.New("unimplemented")
}

func (im IBCModule) OnChanOpenAck(
	ctx sdk.Context,
	portID,
	channelID string,
	counterpartyChannelID string,
	counterpartyVersion string,
) error {
	return errors.New("unimplemented")
}

func (im IBCModule) OnChanOpenConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return errors.New("unimplemented")
}

func (im IBCModule) OnChanCloseInit(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return errors.New("unimplemented")
}

func (im IBCModule) OnChanCloseConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return errors.New("unimplemented")
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
	return errors.New("unimplemented")
}

func (im IBCModule) OnTimeoutPacket(
	ctx sdk.Context,
	packet ibcchanneltypes.Packet,
	relayer sdk.AccAddress,
) error {
	return errors.New("unimplemented")
}
