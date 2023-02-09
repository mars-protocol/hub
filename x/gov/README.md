# Custom Gov

The `customgov` module is wrapper around Cosmos SDK's vanilla `gov` module, inheriting most of its functionalities, with two changes:

- A custom vote tallying logic. Namely, tokens locked in the [vesting contract](https://github.com/mars-protocol/hub-periphery/tree/main/contracts/vesting) count towards one's governance voting power.
- Type check proposal and vote metadata.

## Tallying

Let's illustrate this by an example.

Consider a blockchain with only one validator, as well as two users, Alice and Bob, who have the following amounts of tokens:

| user  | staked | vesting |
| ----- | ------ | ------- |
| Alice | 30     | 21      |
| Bob   | 49     | 0       |

Assume the validator votes YES on a proposal, while neither Alice or Bob votes. In this case, the validator will vote on behalf of Alice and Bob's staked tokens. The vote will pass with 30 + 49 = 79 tokens voting YES and the rest 21 tokens not voting.

If Alice votes NO, this overrides the validator's voting. The vote will be defeated by 49 tokens voting YES vs 51 tokens voting NO.

### Note

Currently, the module assumes the vesting contract is the first contract to be deployed on the chain, i.e. having the code ID of 1 and instance ID of 1. The module uses this info [to derive the contract's address](https://github.com/mars-protocol/hub/blob/2d233fe074b008c49cf26362e1446d888fc81ca0/custom/gov/keeper/tally.go#L12-L15). Developers must make sure this is the case in the chain's genesis state.

Why not make it a configurable parameter? Because doing so involves modifying gov module's `Params` type definition which breaks a bunch of things, which we prefer not to.

## Metadata

From Cosmos SDK v0.46, governance proposals no longer have a "title" and a "description", but instead a "metadata" which can be an arbitrary string. According to [the docs](https://docs.cosmos.network/main/modules/gov#proposal-3), the recommended way to provide the metadata is to store it off-chain, and only upload an IPFS hash on-chain. Therefore, the vanilla gov module:

- Has a default 255 character limit for the metadata string
- Does not enforce a schema of the metadata string

In Mars `customgov`, we want to storage the metadata on-chain. In order for this to work, we increase the length limit to `u64::MAX`, essentially without a limit. Additionally, we implement type checks for the metadata. Specifically,

- For proposal metadata, we assert that it is non-empty and conforms to this schema (defined in TypeScript):

  ```typescript
  type ProposalMetadata = {
    title: string;
    authors?: string[];
    summary: string;
    details?: string;
    proposal_forum_url?: string;
    vote_option_context?: string;
  };
  ```

  We make `title` and `summary` mandatory and the other fields optional, because from sdk 0.47 [proposals will have mandatory title and summary fields](https://github.com/cosmos/cosmos-sdk/blob/v0.47.0-rc1/proto/cosmos/gov/v1/gov.proto#L85-L93). Once Mars Hub upgrades to sdk 0.47, we can make these two fields optional as well.

- For vote metadata, we assert that it is either an empty string (it's ok if a voter doesn't want to provide a rationale for their vote), or if it's not empty, conforms to this schema:

  ```typescript
  type VoteMetadata = {
    justification?: string;
  };
  ```
