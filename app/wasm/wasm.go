package wasm

import (
	wasm "github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"

	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
)

func RegisterCustomPlugins(distrKeeper *distrkeeper.Keeper) []wasm.Option {
	messengerDecoratorOpt := wasmkeeper.WithMessageHandlerDecorator(
		CustomMessageDecorator(distrKeeper),
	)

	queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom: CustomQuerier(&QueryPlugin{}),
	})

	return []wasm.Option{messengerDecoratorOpt, queryPluginOpt}
}
