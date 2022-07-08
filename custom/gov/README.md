# Custom Gov

The `customgov` module is wrapper around Cosmos SDK's vanilla `gov` module, inheriting most of its functionalities, but replacing the vote tallying logic with our custom implementation. Namely, tokens locked in the [vesting contract](https://github.com/mars-protocol/hub-periphery/tree/main/contracts/vesting) count towards own's governance voting power.

## Example

Consider a blockchain with only one validator, as well as two users, Alice and Bob, who have the following amounts of tokens:

| user  | staked | vesting |
| ----- | ------ | ------- |
| Alice | 30     | 21      |
| Bob   | 49     | 0       |

Assume the validator votes YES on a proposal, while neither Alice or Bob votes. In this case, the validator will vote on behalf of Alice and Bob's staked tokens. The vote will pass with 30 + 49 = 79 tokens voting YES and the rest 21 tokens not voting.

If Alice votes NO, this overrides the validator's voting. The vote will be defeated by 49 tokens voting YES vs 51 tokens voting NO.

## Note

Currently, the module assumes the vesting contract is the first contract to be deployed on the chain, i.e. having the code ID of 1 and instance ID of 1. The module uses this info to derive the contract's address. Developers must make sure this is the case in the chain's genesis state.

For future releases, it may be a good idea to make the contract address a configurable parameter.
