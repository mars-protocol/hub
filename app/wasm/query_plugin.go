package wasm

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
)

func CustomQuerier(qp *QueryPlugin) func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		var marsQuery MarsQuery
		if err := json.Unmarshal(request, &marsQuery); err != nil {
			return nil, sdkerrors.Wrapf(err, "invalid x query: %s", request)
		}

		// here, dispatch query request to the appropriate query function

		return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown x query variant"}
	}
}

type QueryPlugin struct {
	// currently we don't have any x query implemented
}
