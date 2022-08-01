# Mars Hub

Mars Hub application-specific blockchain, built on top of [Tendermint](https://github.com/tendermint/tendermint), [Cosmos SDK](https://github.com/cosmos/cosmos-sdk), [IBC](https://github.com/cosmos/ibc-go), and [CosmWasm](https://github.com/CosmWasm/wasmd).

## Installation

Install the lastest version of [Go programming language](https://go.dev/dl/) and configure related environment variables. See [here](https://github.com/st4k3h0us3/workshops/tree/main/how-to-run-a-validator) for a tutorial.

Clone this repository, checkout to the latest tag, the compile the code:

```bash
git clone https://github.com/mars-protocol/hub.git
cd hub
git checkout <tag>
make install
```

A `marsd` executable will be created in the `$GOBIN` directory

## License

TBD
