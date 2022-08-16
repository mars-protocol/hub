# relay

`relay` is the [authentication module](https://ibc.cosmos.network/main/apps/interchain-accounts/auth-modules.html) of Mars Hub's [Interchain Account (ICA)](https://github.com/cosmos/ibc/blob/main/spec/app/ics-027-interchain-accounts/README.md) Controller.

- If a new Outpost is to be deployed, the relay module will register a new ICA on the destination chain, which will act as the owner and admin of the Outpost contracts.

- The relay module comes with two governance proposal types: `ExecuteRemoteContractProposal` and `MigrateRemoteContractProposal`. If one such proposal (e.g. adding new asset to an Outpost, adjusting risk parameters, or migrating contracts) is passed in governance, the relay module will dispatch the appropriate wasm message(s) to the Outpost chain to be executed via ICA.

The name "relay" is inspired by [mass relays](https://masseffect.fandom.com/wiki/Mass_Relay) from the _Mass Effect_ series, which are devices that can transport matter or data at FTL speed across the galaxy.
