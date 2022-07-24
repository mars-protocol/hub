# Incentives

The incentives module manages incentivization for MARS stakers. **Not to be confused** with incentives for lending/borrowing activities, which are managed by a wasm contract deployed at each Outpost.

The release of incentives are defined by **Schedules**. Each incentives schedule consists three (3) parameters:

- `StartTime`
- `EndTime`
- `TotalAmount`

Between the timespan defined by `StartTime` and `EndTime`, coins specified by `TotalAmount` will be released as staking rewards _linearly_, in the `BeginBlocker` of each block. Each validator _who have signed the previous block_ gets a portion of the block reward pro-rata according to their voting power.

A new schedule can be created upon a successful `CreateIncentivesScheduleProposal`. The incentives module will withdraw the coins corresponding to `TotalAmount` from the community pool to its module account. Conversely, an active schedule can be cancelled upon a successful `TerminateIncentivesScheduleProposal`. All coins yet to be distributed will be returned to the community pool.

There can be multiple schedules active at the same time, each identified by a `uint64`. Each schedule can release multiple coins, not limited to the MARS token.
