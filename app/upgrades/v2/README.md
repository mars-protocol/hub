# v2 upgrade

In v2.0.0 upgrade we add the interchain account and envoy modules, without making any change to the existing modules.

## Manual testing

Here are some steps I took to test this upgrade using two local devnets. The objective is to verify:

- ICS-20 channels created pre-upgrade are still active
- interchain account modules are initialized properly
- envoy module works

### marsd

First, check out to `v1.0.0` tag and compile:

```bash
git checkout v1.0.0
make build
mv build/marsd $GOBIN/marsd-v1.0.0
marsd-v1.0.0 version --long
```

Create genesis state:

```bash
marsd-v1.0.0 init testnode --chain-id mars-test-1 --staking-bond-denom umars --home ~/.mars-v2-upgrade-test
marsd-v1.0.0 config keyring-backend test --home ~/.mars-v2-upgrade-test
```

Create an account. We will use this account for everything: validator, relayer, proposal submitter...

```bash
marsd-v1.0.0 keys add test --home ~/.mars-v2-upgrade-test
testAcc=$(marsd-v1.0.0 keys show test --address --home ~/.mars-v2-upgrade-test)
```

Set up genesis:

```bash
marsd-v1.0.0 genesis add-account $testAcc 100000000umars --home ~/.mars-v2-upgrade-test
marsd-v1.0.0 genesis gentx test 5000000umars --home ~/.mars-v2-upgrade-test
marsd-v1.0.0 genesis collect-gentxs --home ~/.mars-v2-upgrade-test
marsd-v1.0.0 genesis add-wasm-message store ~/.mars-v2-upgrade-test/mars_vesting.wasm --run-as $testAcc --home ~/.mars-v2-upgrade-test
marsd-v1.0.0 genesis add-wasm-message instantiate-contract 1 "{\"owner\":\"$testAcc\",\"unlock_schedule\":{\"start_time\":0,\"cliff\":0,\"duration\":1}}" --label mars-vesting --run-as $testAcc --no-admin --home ~/.mars-v2-upgrade-test
```

Edit `genesis.json`, setting gov voting period to `60s`.

### wasmd

We use wasmd v0.31.0 as the counterparty chain. Similar set up:

```bash
wasmd init testnode --chain-id wasm-test-1
wasmd config keyring-backend test
wasmd config node tcp://localhost:36657
wasmd keys add test
wasmTestAcc=$(wasmd keys show test --address)
wasmd add-genesis-account $wasmTestAcc 100000000uwasm
wasmd gentx test 5000000uwasm --chain-id wasm-test-1
wasmd collect-gentxs
```

Edit `~/.wasmd/config/genesis.json`, setting bond denom to `uwasm`.

Edit `~/.wasmd/config/app.toml` and `config.toml`, setting the correct ports.

### hermes

Install hermes:

```bash
cargo install ibc-relayer-cli --bin hermes --locked
```

Save the following file to `~/.hermes/config.toml`:

```toml
[global]
log_level = 'info'

[mode]

[mode.clients]
enabled = true
refresh = true
misbehaviour = false

[mode.connections]
enabled = false

[mode.channels]
enabled = false

[mode.packets]
enabled = true
clear_interval = 100
clear_on_start = true
tx_confirmation = false
auto_register_counterparty_payee = false

[rest]
enabled = false
host = '127.0.0.1'
port = 3000

[telemetry]
enabled = false
host = '127.0.0.1'
port = 3001

[[chains]]
id = 'mars-test-1'
rpc_addr = 'http://127.0.0.1:26657'
grpc_addr = 'http://127.0.0.1:9090'
websocket_addr = 'ws://127.0.0.1:26657/websocket'
rpc_timeout = '10s'
account_prefix = 'mars'
key_name = 'relayer'
store_prefix = 'ibc'
default_gas = 100000
max_gas = 20000000
gas_price = { price = 0, denom = 'umars' }
gas_multiplier = 1.1
max_msg_num = 30
max_tx_size = 2097152
clock_drift = '5s'
max_block_time = '30s'
trusting_period = '14days'
trust_threshold = { numerator = '1', denominator = '3' }
address_type = { derivation = 'cosmos' }

[[chains]]
id = 'wasm-test-1'
rpc_addr = 'http://127.0.0.1:36657'
grpc_addr = 'http://127.0.0.1:9091'
websocket_addr = 'ws://127.0.0.1:36657/websocket'
rpc_timeout = '10s'
account_prefix = 'wasm'
key_name = 'relayer'
store_prefix = 'ibc'
default_gas = 100000
max_gas = 20000000
gas_price = { price = 0, denom = 'umars' }
gas_multiplier = 1.1
max_msg_num = 30
max_tx_size = 2097152
clock_drift = '5s'
max_block_time = '30s'
trusting_period = '14days'
trust_threshold = { numerator = '1', denominator = '3' }
address_type = { derivation = 'cosmos' }
```

Import keys to hermes:

```bash
hermes keys add --chain mars-test-1 --mnemonic-file ~/.hermes/mnemonic-mars.txt
hermes keys add --chain wasm-test-1 --mnemonic-file ~/.hermes/mnemonic-wasm.txt
```

### Run v1.0.0

Start the chains and relayer:

```bash
marsd-v1.0.0 start --home ~/.mars-v2-upgrade-test
wasmd start
hermes start
```

Create a transfer channel. We want to verify this channel is operational before and after the upgrade:

```bash
hermes create channel --a-chain mars-test-1 --b-chain wasm-test-1 --a-port transfer --b-port transfer --new-client-connection
```

Try sending an IBC transfer:

```bash
marsd-v1.0.0 tx ibc-transfer transfer transfer channel-0 $wasmTestAcc 12345umars --from test --gas auto --gas-adjustment 1.4 --home ~/.mars-v2-upgrade-test
wasmd q bank balances $wasmTestAcc
```

Submit a software upgrade proposal and vote:

```bash
cat << EOF >> proposal.json
{
  "messages": [
    {
      "@type": "/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade",
      "authority": "mars10d07y265gmmuvt4z0w9aw880jnsr700j8l2urg",
      "plan": {
        "name": "v2",
        "height": "500"
      }
    }
  ],
  "metadata": "{\"title\":\"Upgrade Mars Hub to v2\",\"summary\":\"n/a\"}",
  "deposit": "10000000umars"
}
EOF
marsd-v1.0.0 tx gov submit-proposal proposal.json --from test --gas auto --gas-adjustment 1.4 --home ~/.mars-v2-upgrade-test
marsd-v1.0.0 tx gov vote 1 yes --from test --gas auto --gas-adjustment 1.4 --home ~/.mars-v2-upgrade-test
marsd-v1.0.0 q gov proposal 1 --home ~/.mars-v2-upgrade-test
marsd-v1.0.0 q upgrade plan --home ~/.mars-v2-upgrade-test
```

### Upgrade to v2.0.0

Marsd should halt at height 500 with the following log msg:

```plain
11:53PM ERR UPGRADE "v2" NEEDED at height: 500
11:53PM ERR CONSENSUS FAILURE!!! err="UPGRADE \"v2\" NEEDED at height: 500"
```

Check out to `v2.0.0` and compile:

```bash
git checkout v2.0.0
make install
marsd version --long
```

Restart marsd and hermes:

```bash
marsd start --home ~/.mars-v2-upgrade-test
hermes start
```

Should see the following log msgs:

```plain
11:54PM INF applying upgrade "v2" at height: 500
11:54PM INF 🚀 executing Mars Hub v2 upgrade 🚀
11:54PM INF initializing interchain account module
11:54PM INF created new capability module=ibc name=ports/icahost
11:54PM INF port binded module=x/ibc/port port=icahost
11:54PM INF claimed capability capability=3 module=icahost name=ports/icahost
11:54PM INF initializing envoy module
```

Make sure new modules are successfully initiated:

```bash
marsd q interchain-accounts controller params --home ~/.mars-v2-upgrade-test
marsd q interchain-accounts host params --home ~/.mars-v2-upgrade-test
marsd q auth module-account envoy --home ~/.mars-v2-upgrade-test
```

Try sending an IBC transfer again, should work:

```bash
marsd-v1.0.0 tx ibc-transfer transfer transfer channel-0 $wasmTestAcc 12345umars --from test --gas auto --gas-adjustment 1.4 --home ~/.mars-v2-upgrade-test
wasmd q bank balances $wasmTestAcc
```

Attempt to register an ICA on wasm-test-1:

```bash
marsd tx envoy register-account connection-0 --from test --gas auto --gas-adjustment 1.4 --home ~/.mars-v2-upgrade-test
```

In my case hermes didn't pick up the channel handshake events for some reason, so I had to do it manually:

```bash
envoyModAcc="mars1s3fjkvr0yk2c0smyh4esrcyp893atwz0uga6lf"
controllerPort="icacontroller-$envoyModAcc"
hermes tx chan-open-try --dst-chain wasm-test-1 --src-chain mars-test-1 --dst-connection connection-0 --dst-port icahost --src-port $controllerPort --src-channel channel-1
hermes tx chan-open-ack --dst-chain mars-test-1 --src-chain wasm-test-1 --dst-connection connection-0 --dst-port $controllerPort --src-port icahost --dst-channel channel-1 --src-channel channel-1
hermes tx chan-open-confirm --dst-chain wasm-test-1 --src-chain mars-test-1 --dst-connection connection-0 --dst-port icahost --src-port $controllerPort --dst-channel channel-1 --src-channel channel-1
```

Reboot hermes, it should correctly detect the channels:

```plain
2023-03-16T00:03:02.977023Z  INFO ThreadId(01) # Chain: mars-test-1
  - Client: 07-tendermint-0
    * Connection: connection-0
      | State: OPEN
      | Counterparty state: OPEN
      + Channel: channel-0
        | Port: transfer
        | State: OPEN
        | Counterparty: channel-0
      + Channel: channel-1
        | Port: icacontroller-mars1s3fjkvr0yk2c0smyh4esrcyp893atwz0uga6lf
        | State: OPEN
        | Counterparty: channel-1
# Chain: wasm-test-1
  - Client: 07-tendermint-0
    * Connection: connection-0
      | State: OPEN
      | Counterparty state: OPEN
      + Channel: channel-0
        | Port: transfer
        | State: OPEN
        | Counterparty: channel-0
      + Channel: channel-1
        | Port: icahost
        | State: OPEN
        | Counterparty: channel-1
```

Do some queries to make sure the account is properly registered:

```bash
marsd q envoy account connection-0 --home ~/.mars-v2-upgrade-test
marsd q envoy accounts --home ~/.mars-v2-upgrade-test
marsd q interchain-accounts controller interchain-account $envoyModAcc connection-0 --home ~/.mars-v2-upgrade-test
marsd q ibc channel end icacontroller-mars1s3fjkvr0yk2c0smyh4esrcyp893atwz0uga6lf channel-1 --home ~/.mars-v2-upgrade-test
```
