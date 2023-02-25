<!-- markdownlint-configure-file { "no-inline-html": false } -->

# E2E tests for Envoy module

This directory contains end-to-end tests for the Envoy module.

The tests are written in [Typescript][1] to be executable in [Deno runtime][2].

## Executing

This end-to-end test is set up as a [GitHub Actions workflow][3].

There are two ways to run this _workflow_ on code changes.

### CI/CD

Push commits to any GitHub repository will trigger this workflow on GitHub actions.

### Locally

[`act`][4] can be used to test this workflow on a local machine.

```sh
act -j envoy
```

### Manually

It is also possible to run this test manually, outside of the Deno runtime, directing interacting with the node and relayer software. See [`MANUAL.md`][5] for instructions.

## Description

This script spawns up one `marsd` and one `wasmd` with a ibc-transafer channel opened between them and `hermes` relayer, relaying ibc-packets between them.

Then the script performs assertions for the transaction and query API of the Envoy module.

- Account Registration: ICA registration of Envoy module account on the counterparty blockchain. It also asserts the query APIs.
- Send Funds: Send funds from the Envoy module account to its ICA on the counterparty blockchain.
- Submit ICA Transactions: Submit transactions from the Envoy module account to its ICA on the counterparty blockchain.

_Send Funds_ and _Submit ICA Transactions_ can not be submitted directly as a signed transaction. So they are executed through accepted governance proposals.

<sub>
Deno is preferred because, unlike other programming languages, Deno [resolves dependencies][6] on the fly and fetches them if needed. This declutters project space and avoids extra dependency installation steps.
</sub>

[1]: https://www.typescriptlang.org
[2]: https://deno.land
[3]: ../../.github/workflows/e2e.yml
[4]: https://github.com/nektos/act
[5]: ./MANUAL.md
[6]: https://deno.land/manual/examples/manage_dependencies
