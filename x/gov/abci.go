package gov

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	"github.com/mars-protocol/hub/x/gov/keeper"
)

// EndBlocker called at the end of every block, processing proposals
//
// This is pretty much the same as the vanilla gov EndBlocker, except for we
// replace the `Tally` function with our own implementation.
func EndBlocker(ctx sdk.Context, keeper keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(govtypes.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	logger := keeper.Logger(ctx)

	// Delete dead proposals from store and returns theirs deposits.
	// A proposal is dead when it's inactive and didn't get enough deposit on
	// time to get into voting phase.
	keeper.IterateInactiveProposalsQueue(ctx, ctx.BlockHeader().Time, func(proposal govv1.Proposal) bool {
		keeper.DeleteProposal(ctx, proposal.Id)
		keeper.RefundAndDeleteDeposits(ctx, proposal.Id)

		// called when proposal become inactive
		keeper.AfterProposalFailedMinDeposit(ctx, proposal.Id)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				govtypes.EventTypeInactiveProposal,
				sdk.NewAttribute(govtypes.AttributeKeyProposalID, fmt.Sprintf("%d", proposal.Id)),
				sdk.NewAttribute(govtypes.AttributeKeyProposalResult, govtypes.AttributeValueProposalDropped),
			),
		)

		logger.Info(
			"proposal did not meet minimum deposit; deleted",
			"proposal", proposal.Id,
			"min_deposit", sdk.NewCoins(keeper.GetDepositParams(ctx).MinDeposit...).String(),
			"total_deposit", sdk.NewCoins(proposal.TotalDeposit...).String(),
		)

		return false
	})

	// fetch active proposals whose voting periods have ended (are passed the block time)
	keeper.IterateActiveProposalsQueue(ctx, ctx.BlockHeader().Time, func(proposal govv1.Proposal) bool {
		var tagValue, logMsg string

		// IMPORTANT: use our custom implementation of tally logics
		passes, burnDeposits, tallyResults := keeper.Tally(ctx, proposal)

		if burnDeposits {
			keeper.DeleteAndBurnDeposits(ctx, proposal.Id)
		} else {
			keeper.RefundAndDeleteDeposits(ctx, proposal.Id)
		}

		if passes {
			var (
				idx    int
				events sdk.Events
				msg    sdk.Msg
			)

			// attempt to execute all messages within the passed proposal
			// Messages may mutate state thus we use a cached context. If one of
			// the handlers fails, no state mutation is written and the error
			// message is logged.
			cacheCtx, writeCache := ctx.CacheContext()
			messages, err := proposal.GetMsgs()
			if err == nil {
				for idx, msg = range messages {
					handler := keeper.Router().Handler(msg)

					var res *sdk.Result
					res, err = handler(cacheCtx, msg)
					if err != nil {
						break
					}

					events = append(events, res.GetEvents()...)
				}
			}

			// `err == nil` when all handlers passed.
			// Or else, `idx` and `err` are populated with the msg index and error.
			if err == nil {
				proposal.Status = govv1.StatusPassed
				tagValue = govtypes.AttributeValueProposalPassed
				logMsg = "passed"

				// write state to the underlying multi-store
				writeCache()

				// propagate the msg events to the current context
				ctx.EventManager().EmitEvents(events)
			} else {
				proposal.Status = govv1.StatusFailed
				tagValue = govtypes.AttributeValueProposalFailed
				logMsg = fmt.Sprintf("passed, but msg %d (%s) failed on execution: %s", idx, sdk.MsgTypeURL(msg), err)
			}
		} else {
			proposal.Status = govv1.StatusRejected
			tagValue = govtypes.AttributeValueProposalRejected
			logMsg = "rejected"
		}

		proposal.FinalTallyResult = &tallyResults

		keeper.SetProposal(ctx, proposal)
		keeper.RemoveFromActiveProposalQueue(ctx, proposal.Id, *proposal.VotingEndTime)

		// when proposal become active
		keeper.AfterProposalVotingPeriodEnded(ctx, proposal.Id)

		logger.Info(
			"proposal tallied",
			"proposal", proposal.Id,
			"results", logMsg,
		)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				govtypes.EventTypeActiveProposal,
				sdk.NewAttribute(govtypes.AttributeKeyProposalID, fmt.Sprintf("%d", proposal.Id)),
				sdk.NewAttribute(govtypes.AttributeKeyProposalResult, tagValue),
			),
		)

		return false
	})
}
