# Mars Hub

Mars Hub application-specific blockchain


## Example usage scenario
*contract on osmo <--------> mars hub <----------> contact on juno
 

## Development note:

"custom" is a folder for customized bits of cosmos sdk modules.  Instead of forking or even using the whole module, mars 
## Installation

Install the lastest version of [Go programming language](https://go.dev/dl/) and configure related environment variables. See [here](https://github.com/st4k3h0us3/workshops/tree/main/how-to-run-a-validator) for a tutorial.

Clone this repository, checkout to the latest tag, the compile the code:

```bash
git clone https://github.com/mars-protocol/hub.git
cd hub
git checkout <tag>
make install
```

To use [the experimental RocksDB backend](https://github.com/tendermint/tm-db/pull/237):

```bash
make MARS_BUILD_OPTIONS='rocksdb' install
```

A `marsd` executable will be created in the `$GOBIN` directory

## License

TBD
