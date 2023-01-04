# Custom Gov

The `customgov` module is wrapper around Cosmos SDK's vanilla `gov` module, inheriting most of its functionalities, with two changes:

- Replacing the vote tallying logic with our custom implementation. Namely, tokens locked in the [vesting contract](https://github.com/mars-protocol/hub-periphery/tree/main/contracts/vesting) count towards one's governance voting power.
- Type check proposal and vote metadata.

## Vote tallying logic

Let's illustrate this by an example.

Consider a blockchain with only one validator, as well as two users, Alice and Bob, who have the following amounts of tokens:

| user  | staked | vesting |
| ----- | ------ | ------- |
| Alice | 30     | 21      |
| Bob   | 49     | 0       |

Assume the validator votes YES on a proposal, while neither Alice or Bob votes. In this case, the validator will vote on behalf of Alice and Bob's staked tokens. The vote will pass with 30 + 49 = 79 tokens voting YES and the rest 21 tokens not voting.

If Alice votes NO, this overrides the validator's voting. The vote will be defeated by 49 tokens voting YES vs 51 tokens voting NO.

### A note

Currently, the module assumes the vesting contract is the first contract to be deployed on the chain, i.e. having the code ID of 1 and instance ID of 1. The module uses this info [to derive the contract's address](https://github.com/mars-protocol/hub/blob/2d233fe074b008c49cf26362e1446d888fc81ca0/custom/gov/keeper/tally.go#L12-L15). Developers must make sure this is the case in the chain's genesis state.

Why not make it a configurable parameter? Because doing so involves modifying gov module's `Params` type definition which breaks a bunch of things, which we prefer not to.

## Proposal metadata

From cosmos-sdk 0.46, governance proposals no longer have a "title" and a "description", but instead a "metadata" which can be an arbitrary string, and by default has a length limit of 255 characters.

This is a big mistake! For two reasons:

- The 255 character limit is there because the module assumes the proposals will be stored off-chain, and only an IPFS hash will be uploaded on-chain. In my opinion, proposals are an important part of a blockchain's history and should be persisted on-chain. In our custom gov module we set the length limit to `u64::MAX`.
- The module doesn't enforce a type/format/schema of the metadata. This creates problems for webapps, wallets, and block explorers because they won't know how to parse the metadata. At Mars we implement a type-check of the metadata as a part of the handling of `MsgSubmitProposal`.

We require all proposals to come with a non-empty metadata string with this format (defined in TypeScript):

```typescript
type ProposalMetadata = {
  title: string;
  authors: string[];
  summary?: string;
  details: string;
  proposal_forum_url?: string;
  vote_option_context?: string;
};
```

Note that this differs from [what the SDK docs recommands](https://docs.cosmos.network/main/modules/gov#proposal-3) in that `authors` is an array of string instead of a single string, and that a few fields are made optional. We will make a PR to the SDK repo with this change. Hopefully it'll be accepted.

For votes, the metadata can be empty, but if it's non-empty, it must conform to this format:

```typescript
type VoteMetadata = {
  justification?: string;
};
```
