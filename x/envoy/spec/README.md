# TLA+ Specification for Mars Hub's Envoy module

Requires [Apalache model checker](https://apalache.informal.systems).

## Invariants

| ID             | Invariant         | Description                                                        |
| -------------- | ----------------- | ------------------------------------------------------------------ |
| `invariant-CS` | `ConstantSupply`  | Token supply on the chain never changes.                           |
| `invariant-PB` | `PositiveBalance` | Balances are always non-negative.                                  |
| `invariant-VA` | `ValidAuthority`  | SendFunds and SendMessages have the correct authority.             |
| `invariant-IE` | `ICAExists`       | SendFunds and SendMessages are submitted to an existing ICAccount. |
| `invariant-PL` | `NoPacketLoss`    | IBC packets of SendFunds and SendMessages are in IBC queue.        |

`InvAll` is a conjunction of all of them. So to check all of them together,

```sh
apalache-mc check --inv=InvAll envoy.tla
```

## Examples

| ID           | Property     | Invariant      | Description                                                       |
| ------------ | ------------ | -------------- | ----------------------------------------------------------------- |
| `example-AS` | `AllSuccess` | `ExAllSuccess` | Example trace has minimum 5 states and all the actions succeeded. |

### Views

| ID        | Operator     | Description                              |
| --------- | ------------ | ---------------------------------------- |
| `view-AT` | `ActionType` | Projects a state to action message type. |

To generate an example satisfying `AllSuccess`,

```sh
apalache-mc check --inv=ExAllSuccess --max-error=3 --view=ActionType --out-dir=runs envoy.tla
```

The counter-examples would be present at `./runs/envoy.tla/<APALACHE_RUN_ID>/violation*.tla`.
