<!-- markdownlint-configure-file { "no-inline-html": false } -->

# E2E tests for Envoy module

This document describes how I do a manual end-to-end test with two local devnets. I use [marsd][1] (equipped with the envoy module) as the controller chain and [wasmd][2] as the host chain. Use [hermes][3] to relay messages between the two.

## Info

### Ports

We will run marsd and wasmd both on localhost. They will be configured to run using the following ports. Make sure there are no other processes on your computer that occupy the same ports:

|       | marsd | wasmd |
| ----- | ----- | ----- |
| abci  | 26658 | 36658 |
| rpc   | 26657 | 36657 |
| p2p   | 26656 | 36656 |
| pprof | 6060  | 6061  |
| grpc  | 9090  | 9091  |

### Accounts

Each chain comes with the following accounts. They use the `test` keyring backend so don't use them in production:

| Name  | Address                                                                                            | Purpose   |
| ----- | -------------------------------------------------------------------------------------------------- | --------- |
| test1 | `mars1s7mkj5j9jejlqx53dhjx82ljhp4lh4hc2l09nd` <br /> `wasm1s7mkj5j9jejlqx53dhjx82ljhp4lh4hca78f0a` | validator |
| test2 | `mars1hyjtwtdnleadyyp73c3nnz4zs5yfnurue3mwjx` <br /> `wasm1hyjtwtdnleadyyp73c3nnz4zs5yfnuruwsnzwk` | relayer   |
| test3 | `mars14whu3e6dujyhh424nc7tes3q97r9zzdddlhplq` <br /> `wasm14whu3e6dujyhh424nc7tes3q97r9zzdd67ldrs` | user      |

There seed phrases are as follows

```plain
test1
remove pyramid muffin alcohol quit tip situate feed solve urban attend clinic pelican tribe novel task need blanket bamboo join sudden left tunnel faint

test2
panther tree salute panther long cave green build arrow glad champion venture foam magnet tongue grace fun day mixed taxi island emerge kangaroo ribbon

test3
truck inhale decline orphan phrase arena then ahead harsh fortune clinic reveal tomato sick child glow laptop current future task another street once baby
```

### Host

For wasmd there is a contract deployed during genesis at address

```plain
wasm14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s0phg4d
```

This is a barebone contract taken from my [CosmWasm template][4], which uses the [cw-ownable][5] library for ownership management. The owner is initially set to `test3`.

### Objective

The test is considered successful if we can:

- Register an interchain account owned by marsd envoy module on wasmd
- Send coins from marsd community pool to the ICA
- Have the ICA claim the contract's ownership

## Preparation

Download and install marsd (use the main branch):

```bash
git clone https://github.com/mars-protocol/hub/v2.git
cd hub
make install
```

Download and install wasmd (for this test we use v0.30.0):

```bash
git clone https://github.com/CosmWasm/wasmd.git
cd wasmd
git checkout v0.30.0
make install
```

Install hermes:

```bash
cargo install ibc-relayer-cli --bin hermes --locked
```

Extract config folders (`.mars`, `.wasmd`, `.hermes`) to your home directory. Note that this overwrites your local ones so backup first!

```bash
tar xvzf configs.tar.gz -C ~
```

Start the chains

```bash
marsd start
wasmd start
```

Create an ICS-20 transfer channel between the two chains:

```bash
hermes create channel \
  --a-chain mars-dev-1 --a-port transfer \
  --b-chain wasm-dev-1 --b-port transfer \
  --new-client-connection
```

This should create client `07-tendermint-0`, connection `connection-0`, and channel `channel-0` on both chains.

Start hermes:

```bash
hermes start
```

## Register interchain account

```bash
marsd tx envoy register-account connection-0 \
  --from test3 \
  --gas auto \
  --gas-adjustment 1.4
```

Marsd CLI should output the following logging info:

```plain
10:07PM INF created new capability module=ibc name=ports/icacontroller-mars1fr6zyc9ggjx2575u5xhjvv9qsmdy93y0n9670d
10:07PM INF port binded module=x/ibc/port port=icacontroller-mars1fr6zyc9ggjx2575u5xhjvv9qsmdy93y0n9670d
10:07PM INF claimed capability capability=4 module=icacontroller name=ports/icacontroller-mars1fr6zyc9ggjx2575u5xhjvv9qsmdy93y0n9670d
10:07PM INF created new capability module=ibc name=capabilities/ports/icacontroller-mars1fr6zyc9ggjx2575u5xhjvv9qsmdy93y0n9670d/channels/channel-1
10:07PM INF claimed capability capability=5 module=icacontroller name=capabilities/ports/icacontroller-mars1fr6zyc9ggjx2575u5xhjvv9qsmdy93y0n9670d/channels/channel-1
10:07PM INF channel state updated channel-id=channel-1 module=x/ibc/channel new-state=INIT port-id=icacontroller-mars1fr6zyc9ggjx2575u5xhjvv9qsmdy93y0n9670d previous-state=NONE
10:07PM INF initiated interchain account channel handshake connectionID=connection-0 module=x/envoy
```

Check if the interchain account has been successfully registered:

```bash
$ marsd q envoy accounts

accounts:
- address: wasm1jwdap5t78w4na2vmdcaszysqcqkhy0suh4nsg0lce70j7g50f2uqkrnyz4
  controller:
    channel_id: channel-1
    client_id: 07-tendermint-0
    connection_id: connection-0
    port_id: icacontroller-mars1fr6zyc9ggjx2575u5xhjvv9qsmdy93y0n9670d
  host:
    channel_id: channel-1
    client_id: 07-tendermint-0
    connection_id: connection-0
    port_id: icahost
```

```bash
$ address=$(marsd q envoy account connection-0 --output json | jq -r '.account.address')
$ wasmd q auth account $address

'@type': /ibc.applications.interchain_accounts.v1.InterchainAccount
account_owner: icacontroller-mars1fr6zyc9ggjx2575u5xhjvv9qsmdy93y0n9670d
base_account:
  account_number: "10"
  address: wasm1jwdap5t78w4na2vmdcaszysqcqkhy0suh4nsg0lce70j7g50f2uqkrnyz4
  pub_key: null
  sequence: "0"
```

## Send funds

Submit and vote on the proposal:

```bash
marsd tx gov submit-proposal ./send_funds.json \
  --from test3 \
  --gas auto \
  --gas-adjustment 1.4
```

```bash
marsd tx gov vote 1 yes \
  --from test1 \
  --gas auto \
  --gas-adjustment 1.4
```

```bash
marsd q gov proposal 1
```

Once the proposal is passed, you should see the following CLI logging info from marsd, and hermes should start relaying the ICS-20 packets:

```plain
10:09PM INF packet sent dst_channel=channel-0 dst_port=transfer module=x/ibc/channel sequence=1 src_channel=channel-0 src_port=transfer
10:09PM INF IBC fungible token transfer amount=42069 module=x/ibc-transfer receiver=wasm1jwdap5t78w4na2vmdcaszysqcqkhy0suh4nsg0lce70j7g50f2uqkrnyz4 sender=mars1fr6zyc9ggjx2575u5xhjvv9qsmdy93y0n9670d token=uastro
10:09PM INF packet sent dst_channel=channel-0 dst_port=transfer module=x/ibc/channel sequence=2 src_channel=channel-0 src_port=transfer
10:09PM INF IBC fungible token transfer amount=69420 module=x/ibc-transfer receiver=wasm1jwdap5t78w4na2vmdcaszysqcqkhy0suh4nsg0lce70j7g50f2uqkrnyz4 sender=mars1fr6zyc9ggjx2575u5xhjvv9qsmdy93y0n9670d token=umars
10:09PM INF initiated ICS-20 transfer(s) to interchain account amount=42069uastro,69420umars channelID=channel-0 connectionID=connection-0 module=x/envoy
10:09PM INF proposal tallied module=x/gov proposal=1 results=passed
```

Check if the coins have been locked in the transfer module escrow account:

```bash
$ marsd q bank balances $(marsd q ibc-transfer escrow-address transfer channel-0)

balances:
- amount: "42069"
  denom: uastro
- amount: "69420"
  denom: umars
```

Check if the interchain account has received the coins:

```bash
$ wasmd q bank balances $address

balances:
- amount: "42069"
  denom: ibc/20367B8DC0876A0803E1835D8FEE18C7A9ED58DF1EFC3563340DA4628ECB1F6D
- amount: "69420"
  denom: ibc/51A3E9883A23CE8C8AAC802E6FF4A46F8ECDDD361A254B3FB9D71A0EC3FB6A76
```

## Send messages

The contract's current owner proposes to transfer ownership to the ICA:

```bash
wasmd tx wasm execute wasm14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s0phg4d '{"update_ownership":{"transfer_ownership":{"new_owner":"'$address'"}}}' \
  --from test3 \
  --gas auto \
  --gas-adjustment 1.4
```

```bash
$ wasmd q wasm contract-state smart wasm14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s0phg4d '{"ownership":{}}'

data:
  owner: wasm14whu3e6dujyhh424nc7tes3q97r9zzdd67ldrs
  pending_expiry: null
  pending_owner: wasm1jwdap5t78w4na2vmdcaszysqcqkhy0suh4nsg0lce70j7g50f2uqkrnyz4
```

Submit and vote on the proposal:

```bash
marsd tx gov submit-proposal ./send_messages.json \
  --from test3 \
  --gas auto \
  --gas-adjustment 1.4
```

```bash
marsd tx gov vote 2 yes \
  --from test1 \
  --gas auto \
  --gas-adjustment 1.4
```

```bash
marsd q gov proposal 2
```

Once the proposal is passed, you should see the following CLI logging info from marsd, and hermes should start relaying the ICS-27 packets:

```plain
10:11PM INF packet sent dst_channel=channel-1 dst_port=icahost module=x/ibc/channel sequence=1 src_channel=channel-1 src_port=icacontroller-mars1fr6zyc9ggjx2575u5xhjvv9qsmdy93y0n9670d
10:11PM INF initiated ICS-27 tx execution with interchain account connectionID=connection-0 module=x/envoy numMsgs=1
10:11PM INF proposal tallied module=x/gov proposal=2 results=passed
```

Query the contract again the see if the ownership info has been correctly updated:

```bash
$ wasmd q wasm contract-state smart wasm14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s0phg4d '{"ownership":{}}'

data:
  owner: wasm1jwdap5t78w4na2vmdcaszysqcqkhy0suh4nsg0lce70j7g50f2uqkrnyz4
  pending_expiry: null
  pending_owner: null
```

[1]: https://github.com/mars-protocol/hub/v2
[2]: https://github.com/CosmWasm/wasmd
[3]: https://github.com/informalsystems/hermes.git
[4]: https://github.com/steak-enjoyers/cw-template
[5]: https://github.com/steak-enjoyers/cw-plus-plus/tree/main/packages/ownable
