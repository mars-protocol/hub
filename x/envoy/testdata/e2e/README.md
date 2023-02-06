# E2E tests for Envoy module

This directory contains end-to-end tests for the Envoy module.

This end-to-end test is set up to run on
[GitHub workflow CI](../../../../.github/workflows/e2e.yml).

The tests are written in [Typescript](https://www.typescriptlang.org) to be
executable in [Deno runtime](https://deno.land).

## Locally

One may use [`act`](https://github.com/nektos/act) to test this workflow
locally.

Alternatively, one also may run the tests locally using included Dockerfile.

```
docker build . -t envoy-docker
docker run --rm -it envoy-docker deno test -A envoy_test.ts
```

These steps do not require Deno to be installed locally.

> **Note** Docker build takes changes from the MarsHub GitHub repository.

## Description

This script spawns up one `marsd` and one `wasmd` with a ibc-transafer channel
opened between them and `hermes` relayer, relaying ibc-packets between them.

Then the script performs assertions for the transaction and query API of the
Envoy module.

- Account Registration: ICA registration of Envoy module account on the
  counterparty blockchain. It also asserts the query APIs.
- Send Funds: Send funds from the Envoy module account to its ICA on the
  counterparty blockchain.
- Submit ICA Transactions: Submit transactions from the Envoy module account to
  its ICA on the counterparty blockchain.

_Send Funds_ and _Submit ICA Transactions_ can not be submitted directly as a
signed transaction. So they are executed through accepted governance proposals.

<sub>Deno is preferred because, unlike other programming languages, Deno
[resolves dependencies](https://deno.land/manual/examples/manage_dependencies)
on the fly and fetches them if needed. This declutters project space and avoids
extra dependency installation steps.</sub>
