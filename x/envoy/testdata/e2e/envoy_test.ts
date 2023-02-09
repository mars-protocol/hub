import {
  assert,
  assertEquals,
  assertExists,
  assertNotEquals,
} from "https://deno.land/std@0.176.0/testing/asserts.ts";
import {
  afterAll,
  beforeAll,
  describe,
  it,
} from "https://deno.land/std@0.176.0/testing/bdd.ts";
import { retry } from "https://deno.land/std@0.176.0/async/mod.ts";

const VALIDATOR_WALLET = "test1";
const _RELAYER_WALLET = "test2";
const USER_WALLET = "test3";
const SEND_FUNDS_JSON = "../send_funds.json";
const SEND_MSGS_JSON = "../send_messages.json";
const CW_ADDR =
  "wasm14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s0phg4d";

const RETRY_OPTION = {
  multiplier: 1,
  maxTimeout: 90000,
  maxAttempts: 50,
  minTimeout: 2000,
};

const RUNNING_PROCS: Deno.Process[] = [];

async function exec(
  cmd: string[],
  stdin?: string,
  is_detach?: boolean
): Promise<string> {
  const p = Deno.run({
    cmd: cmd,
    stdin: typeof stdin === "undefined" ? "inherit" : "piped",
    stdout: "piped",
  });

  if (typeof p.stdin !== "undefined") {
    await p.stdin!.write(new TextEncoder().encode(stdin));
    await p.stdin!.close();
  }

  let output = "";

  if (typeof is_detach !== "undefined" && is_detach) {
    RUNNING_PROCS.push(p);
  } else {
    try {
      assert((await p.status()).success);
    } finally {
      output = new TextDecoder("utf-8").decode(await p.output());
      await p.close();
    }
  }

  return output;
}

beforeAll(async () => {
  // start marsd
  await exec(["marsd", "start"], undefined, true);

  // start wasmd
  await exec(["wasmd", "start"], undefined, true);

  // wait for minimum 3 blocks
  await retry(async () => {
    await exec(["marsd", "q", "block", "3"]);
  }, RETRY_OPTION);

  // wait for minimum 3 blocks
  await retry(async () => {
    await exec(["wasmd", "q", "block", "3"]);
  }, RETRY_OPTION);

  // create ibc client
  await exec([
    "hermes",
    "create",
    "channel",
    "--a-chain=mars-dev-1",
    "--a-port=transfer",
    "--b-chain=wasm-dev-1",
    "--b-port=transfer",
    "--new-client-connection",
    "--yes",
  ]);

  // start hermes
  await exec(["hermes", "start"], undefined, true);
});

afterAll(async () => {
  while (RUNNING_PROCS.length > 0) {
    const p = RUNNING_PROCS.pop()!;
    p.kill("SIGTERM");
    p.stdout!.close();
    await p.close();
  }
});

describe("e2e tests for envoy module", async () => {
  let interchainAccAddr: string;

  await it("register interchain-account of envoy module", async () => {
    // tx to register interchain-account of envoy module
    await exec([
      "marsd",
      "tx",
      "envoy",
      "register-account",
      "connection-0",
      `--from=${USER_WALLET}`,
      "--gas=auto",
      "--gas-adjustment=1.4",
      "--yes",
    ]);

    // wait until the interchain-account is registered
    await retry(async () => {
      const data = JSON.parse(
        await exec(["marsd", "q", "envoy", "accounts", "--output=json"])
      );

      assert(data.accounts.length > 0);
    }, RETRY_OPTION);

    const icas = JSON.parse(
      await exec(["marsd", "q", "envoy", "accounts", "--output=json"])
    ).accounts;

    // interchain-account is registered
    assert(icas.length > 0);

    interchainAccAddr = JSON.parse(
      await exec([
        "marsd",
        "q",
        "envoy",
        "account",
        "connection-0",
        "--output=json",
      ])
    ).account.address;

    assertEquals(interchainAccAddr, icas[0].address);

    console.log({ ICAAddr: interchainAccAddr });
  });

  await it("send funds to interchain-account", async () => {
    // gov proposal to send funds to envoy interchain-account
    await exec([
      "marsd",
      "tx",
      "gov",
      "submit-proposal",
      SEND_FUNDS_JSON,
      `--from=${USER_WALLET}`,
      "--gas=auto",
      "--gas-adjustment=1.4",
      "--yes",
    ]);

    // wait until the proposal is submitted
    await retry(async () => {
      const data = JSON.parse(
        await exec(["marsd", "q", "gov", "proposals", "--output=json"])
      );

      assert(data.proposals.length > 0);
    }, RETRY_OPTION);

    // tx to vote for the proposal
    await exec([
      "marsd",
      "tx",
      "gov",
      "vote",
      "1",
      "yes",
      `--from=${VALIDATOR_WALLET}`,
      "--gas=auto",
      "--gas-adjustment=1.4",
      "--yes",
    ]);

    // wait until the proposal is passed
    await retry(async () => {
      const data = JSON.parse(
        await exec(["marsd", "q", "gov", "proposal", "1", "--output=json"])
      );

      assertEquals(data.status, "PROPOSAL_STATUS_PASSED");
    }, RETRY_OPTION);

    // wait until the envoy fund transfer is processed
    await retry(async () => {
      const data = JSON.parse(
        await exec([
          "wasmd",
          "q",
          "bank",
          "balances",
          interchainAccAddr,
          "--output=json",
        ])
      );

      assert(data.balances.length > 0);
    }, RETRY_OPTION);

    const escrowAddress = (
      await exec([
        "marsd",
        "q",
        "ibc-transfer",
        "escrow-address",
        "transfer",
        "channel-0",
      ])
    ).trim();

    interface Coin {
      denom: string;
      balance: string;
    }

    const escrowBalances: Coin[] = JSON.parse(
      await exec([
        "marsd",
        "q",
        "bank",
        "balances",
        escrowAddress,
        "--output=json",
      ])
    ).balances;

    // escrow balance is non-empty
    assert(escrowBalances.length > 0);

    const sentFunds: Coin[] = JSON.parse(
      await Deno.readTextFile(SEND_FUNDS_JSON)
    ).messages[0].amount;

    const wasmdBalances: Coin[] = JSON.parse(
      await exec([
        "wasmd",
        "q",
        "bank",
        "balances",
        interchainAccAddr,
        "--output=json",
      ])
    ).balances;

    // escrow balance matches with ibc balance
    assertEquals(sentFunds.length, wasmdBalances.length);

    // sent funds match with ibc balance
    assert(
      await wasmdBalances.reduce(async (val, wb) => {
        const denomHash = wb.denom.split("/", 2)[1];

        const denomTrace = JSON.parse(
          await exec([
            "wasmd",
            "q",
            "ibc-transfer",
            "denom-trace",
            denomHash,
            "--output=json",
          ])
        ).denom_trace;

        return (
          (await val) &&
          sentFunds.some((sb) => {
            return (
              sb.denom == denomTrace.base_denom && sb.balance == wb.balance
            );
          })
        );
      }, Promise.resolve(true))
    );
  });

  await it("submit transactions to interchain-account", async () => {
    const wasmPayloadTx = {
      update_ownership: {
        transfer_ownership: {
          new_owner: interchainAccAddr,
        },
      },
    };

    const wasmPayloadQ = { ownership: {} };

    // tx to transfer ownership of a smart contract
    await exec([
      "wasmd",
      "tx",
      "wasm",
      "execute",
      CW_ADDR,
      JSON.stringify(wasmPayloadTx),
      `--from=${USER_WALLET}`,
      "--gas=auto",
      "--gas-adjustment=1.4",
      "--yes",
    ]);

    // wait until the transaction is processed
    await retry(async () => {
      const data = JSON.parse(
        await exec([
          "wasmd",
          "q",
          "wasm",
          "contract-state",
          "smart",
          CW_ADDR,
          JSON.stringify(wasmPayloadQ),
          "--output=json",
        ])
      );

      assertNotEquals(data.data.pending_owner, null);
    }, RETRY_OPTION);

    const prev = JSON.parse(
      await exec([
        "wasmd",
        "q",
        "wasm",
        "contract-state",
        "smart",
        CW_ADDR,
        JSON.stringify(wasmPayloadQ),
        "--output=json",
      ])
    ).data;

    // the ownership is pending for envoy interchain-account
    assertEquals(prev.pending_owner, interchainAccAddr);

    // gov proposal to submit transactions to interchain-account
    await exec([
      "marsd",
      "tx",
      "gov",
      "submit-proposal",
      SEND_MSGS_JSON,
      `--from=${USER_WALLET}`,
      "--gas=auto",
      "--gas-adjustment=1.4",
      "--yes",
    ]);

    // wait until the proposal is processed
    await retry(async () => {
      const data = JSON.parse(
        await exec(["marsd", "q", "gov", "proposals", "--output=json"])
      );

      assert(data.proposals.length > 1);
    }, RETRY_OPTION);

    // tx to submit vote
    await exec([
      "marsd",
      "tx",
      "gov",
      "vote",
      "2",
      "yes",
      `--from=${VALIDATOR_WALLET}`,
      "--gas=auto",
      "--gas-adjustment=1.4",
      "--yes",
    ]);

    // wait until the proposal is passed
    await retry(async () => {
      const data = JSON.parse(
        await exec(["marsd", "q", "gov", "proposal", "2", "--output=json"])
      );

      assertEquals(data.status, "PROPOSAL_STATUS_PASSED");
    }, RETRY_OPTION);

    // wait until the envoy transaction is processed in host chain
    await retry(async () => {
      const data = JSON.parse(
        await exec([
          "wasmd",
          "q",
          "wasm",
          "contract-state",
          "smart",
          CW_ADDR,
          JSON.stringify(wasmPayloadQ),
          "--output=json",
        ])
      );

      assertEquals(data.data.pending_owner, null);
    }, RETRY_OPTION);

    const curr = JSON.parse(
      await exec([
        "wasmd",
        "q",
        "wasm",
        "contract-state",
        "smart",
        CW_ADDR,
        JSON.stringify(wasmPayloadQ),
        "--output=json",
      ])
    ).data;

    // previous pending owner is the current owner
    assertEquals(curr.owner, prev.pending_owner);
    // current owner is envoy interchain-account
    assertEquals(curr.owner, interchainAccAddr);
    // current prending owner is null
    assertEquals(curr.pending_owner, null);
  });
});
