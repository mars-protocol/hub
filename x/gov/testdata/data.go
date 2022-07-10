package testdata

import _ "embed"

// VestingWasm is bytecode of the [`mars-vesting`](https://github.com/mars-protocol/hub-periphery/tree/main/contracts/vesting) contract
//go:embed vesting.wasm
var VestingWasm []byte
