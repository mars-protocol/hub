package wasm

import (
	"encoding/json"

	"cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
)

func CustomQuerier(qp *QueryPlugin) func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		var marsQuery MarsQuery
		if err := json.Unmarshal(request, &marsQuery); err != nil {
			return nil, errors.Wrapf(err, "invalid custom query: %s", request)
		}

		// here, dispatch query request to the appropriate query function

		return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown custom query variant"}
	}
}

type QueryPlugin struct {
	// currently we don't have any custom query implemented
}
