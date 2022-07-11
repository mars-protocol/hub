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
// This is pretty much the same as the vanilla gov EndBlocker, except for we replace the `Tally`
// function with our own implementation
func EndBlocker(ctx sdk.Context, keeper keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(govtypes.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	logger := keeper.Logger(ctx)

	// delete inactive proposal from store and its deposits
	keeper.IterateInactiveProposalsQueue(ctx, ctx.BlockHeader().Time, func(proposal govv1.Proposal) bool {
		keeper.DeleteProposal(ctx, proposal.Id)
		keeper.DeleteAndBurnDeposits(ctx, proposal.Id)

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
			"title", proposal.GetTitle(),
			"min_deposit", keeper.GetDepositParams(ctx).MinDeposit.String(),
			"total_deposit", proposal.TotalDeposit.String(),
		)

		return false
	})

	// fetch active proposals whose voting periods have ended (are passed the block time)
	keeper.IterateActiveProposalsQueue(ctx, ctx.BlockHeader().Time, func(proposal govv1.Proposal) bool {
		var tagValue, logMsg string

		passes, burnDeposits, tallyResults := keeper.Tally(ctx, proposal) // our custom implementation of tally logics

		if burnDeposits {
			keeper.DeleteAndBurnDeposits(ctx, proposal.Id)
		} else {
			keeper.RefundAndDeleteDeposits(ctx, proposal.Id)
		}

		if passes {
			handler := keeper.Router().Handler(proposal.Id())
			cacheCtx, writeCache := ctx.CacheContext()

			// The proposal handler may execute state mutating logic depending on the proposal content.
			// If the handler fails, no state mutation is written and the error message is logged.
			err := handler(cacheCtx, proposal.GetContent())
			if err == nil {
				proposal.Status = govv1.StatusPassed
				tagValue = govtypes.AttributeValueProposalPassed
				logMsg = "passed"

				// The cached context is created with a new EventManager. However, since the proposal
				// handler execution was successful, we want to track/keep any events emitted, so we
				// re-emit to "merge" the events into the original Context's EventManager.
				ctx.EventManager().EmitEvents(cacheCtx.EventManager().Events())

				// write state to the underlying multi-store
				writeCache()
			} else {
				proposal.Status = govv1.StatusFailed
				tagValue = govtypes.AttributeValueProposalFailed
				logMsg = fmt.Sprintf("passed, but failed on execution: %s", err)
			}
		} else {
			proposal.Status = govv1.StatusRejected
			tagValue = govtypes.AttributeValueProposalRejected
			logMsg = "rejected"
		}

		proposal.FinalTallyResult = tallyResults

		keeper.SetProposal(ctx, proposal)
		keeper.RemoveFromActiveProposalQueue(ctx, proposal.Id, *proposal.VotingEndTime)

		// when proposal become active
		keeper.AfterProposalVotingPeriodEnded(ctx, proposal.Id)

		logger.Info(
			"proposal tallied",
			"proposal", proposal.Id,
			"title", proposal.GetTitle(),
			"result", logMsg,
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
