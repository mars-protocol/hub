package wasm

import (
	wasm "github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"

	customdistrkeeper "github.com/mars-protocol/hub/x/distribution/keeper"
)

func RegisterCustomPlugins(distrKeeper *customdistrkeeper.Keeper) []wasm.Option {
	messengerDecoratorOpt := wasmkeeper.WithMessageHandlerDecorator(
		CustomMessageDecorator(distrKeeper),
	)

	queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom: CustomQuerier(&QueryPlugin{}),
	})

	return []wasm.Option{messengerDecoratorOpt, queryPluginOpt}
}
