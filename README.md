# Mars Hub

Mars Hub app-chain, built on top of [Tendermint][1], [Cosmos SDK][2], [IBC][3], and [CosmWasm][4].

## Bug bounty

A bug bounty is currently open for Mars Hub and peripheral contracts. See details [here](https://immunefi.com/bounty/mars/).

## Audits

See reports [here](https://github.com/mars-protocol/mars-audits/tree/main/hub).

## Installation

Install the lastest version of [Go programming language][5] and configure related environment variables. See [here][6] for a tutorial.

Clone this repository, checkout to the latest tag, the compile the code:

```bash
git clone https://github.com/mars-protocol/hub.git
cd hub
git checkout <tag>
make install
```

A `marsd` executable will be created in the `$GOBIN` directory.

## License

Contents of this repository are open source under [GNU General Public License v3](./LICENSE) or later.

[1]: https://github.com/cometbft/cometbft
[2]: https://github.com/cosmos/cosmos-sdk
[3]: https://github.com/cosmos/ibc-go
[4]: https://github.com/CosmWasm/wasmd
[5]: https://go.dev/dl/
[6]: https://github.com/larry0x/workshops/tree/main/how-to-run-a-validator
