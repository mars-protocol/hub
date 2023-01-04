package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govv046 "github.com/cosmos/cosmos-sdk/x/gov/migrations/v046"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

//------------------------------------------------------------------------------
// queryServer
//------------------------------------------------------------------------------

type queryServer struct{ k Keeper }

// NewQueryServerImpl creates an implementation of the QueryServer interface for
// the given keeper.
func NewQueryServerImpl(k Keeper) govv1.QueryServer {
	return &queryServer{k}
}

func (qs queryServer) TallyResult(goCtx context.Context, req *govv1.QueryTallyResultRequest) (*govv1.QueryTallyResultResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ProposalId == 0 {
		return nil, status.Error(codes.InvalidArgument, "proposal id can not be 0")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	proposal, ok := qs.k.GetProposal(ctx, req.ProposalId)
	if !ok {
		return nil, status.Errorf(codes.NotFound, "proposal %d doesn't exist", req.ProposalId)
	}

	var tallyResult govv1.TallyResult

	switch {
	case proposal.Status == govv1.StatusDepositPeriod:
		tallyResult = govv1.EmptyTallyResult()

	case proposal.Status == govv1.StatusPassed || proposal.Status == govv1.StatusRejected:
		tallyResult = *proposal.FinalTallyResult

	default:
		// proposal is in voting period
		_, _, tallyResult = qs.k.Tally(ctx, proposal) // replace with our custom Tally function
	}

	return &govv1.QueryTallyResultResponse{Tally: &tallyResult}, nil
}

func (qs queryServer) Proposal(goCtx context.Context, req *govv1.QueryProposalRequest) (*govv1.QueryProposalResponse, error) {
	return qs.k.Proposal(goCtx, req)
}

func (qs queryServer) Proposals(goCtx context.Context, req *govv1.QueryProposalsRequest) (*govv1.QueryProposalsResponse, error) {
	return qs.k.Proposals(goCtx, req)
}

func (qs queryServer) Vote(goCtx context.Context, req *govv1.QueryVoteRequest) (*govv1.QueryVoteResponse, error) {
	return qs.k.Vote(goCtx, req)
}

func (qs queryServer) Votes(goCtx context.Context, req *govv1.QueryVotesRequest) (*govv1.QueryVotesResponse, error) {
	return qs.k.Votes(goCtx, req)
}

func (qs queryServer) Params(goCtx context.Context, req *govv1.QueryParamsRequest) (*govv1.QueryParamsResponse, error) {
	return qs.k.Params(goCtx, req)
}

func (qs queryServer) Deposit(goCtx context.Context, req *govv1.QueryDepositRequest) (*govv1.QueryDepositResponse, error) {
	return qs.k.Deposit(goCtx, req)
}

func (qs queryServer) Deposits(goCtx context.Context, req *govv1.QueryDepositsRequest) (*govv1.QueryDepositsResponse, error) {
	return qs.k.Deposits(goCtx, req)
}

//------------------------------------------------------------------------------
// legacyQueryServer
//------------------------------------------------------------------------------

type legacyQueryServer struct{ qs govv1.QueryServer }

// NewQueryServerImpl creates an implementation of the QueryServer interface for the given keeper
func NewLegacyQueryServerImpl(qs govv1.QueryServer) govv1beta1.QueryServer {
	return &legacyQueryServer{qs}
}

func (qs legacyQueryServer) TallyResult(goCtx context.Context, req *govv1beta1.QueryTallyResultRequest) (*govv1beta1.QueryTallyResultResponse, error) {
	resp, err := qs.qs.TallyResult(goCtx, &govv1.QueryTallyResultRequest{
		ProposalId: req.ProposalId,
	})
	if err != nil {
		return nil, err
	}

	tally, err := govv046.ConvertToLegacyTallyResult(resp.Tally)
	if err != nil {
		return nil, err
	}

	return &govv1beta1.QueryTallyResultResponse{Tally: tally}, nil
}

func (qs legacyQueryServer) Proposal(goCtx context.Context, req *govv1beta1.QueryProposalRequest) (*govv1beta1.QueryProposalResponse, error) {
	resp, err := qs.qs.Proposal(goCtx, &govv1.QueryProposalRequest{
		ProposalId: req.ProposalId,
	})
	if err != nil {
		return nil, err
	}

	proposal, err := govv046.ConvertToLegacyProposal(*resp.Proposal)
	if err != nil {
		return nil, err
	}

	return &govv1beta1.QueryProposalResponse{Proposal: proposal}, nil
}

func (qs legacyQueryServer) Proposals(goCtx context.Context, req *govv1beta1.QueryProposalsRequest) (*govv1beta1.QueryProposalsResponse, error) {
	resp, err := qs.qs.Proposals(goCtx, &govv1.QueryProposalsRequest{
		ProposalStatus: govv1.ProposalStatus(req.ProposalStatus),
		Voter:          req.Voter,
		Depositor:      req.Depositor,
		Pagination:     req.Pagination,
	})
	if err != nil {
		return nil, err
	}

	legacyProposals := make([]govv1beta1.Proposal, len(resp.Proposals))
	for idx, proposal := range resp.Proposals {
		legacyProposals[idx], err = govv046.ConvertToLegacyProposal(*proposal)
		if err != nil {
			return nil, err
		}
	}

	return &govv1beta1.QueryProposalsResponse{
		Proposals:  legacyProposals,
		Pagination: resp.Pagination,
	}, nil
}

func (qs legacyQueryServer) Vote(goCtx context.Context, req *govv1beta1.QueryVoteRequest) (*govv1beta1.QueryVoteResponse, error) {
	resp, err := qs.qs.Vote(goCtx, &govv1.QueryVoteRequest{
		ProposalId: req.ProposalId,
		Voter:      req.Voter,
	})
	if err != nil {
		return nil, err
	}

	vote, err := govv046.ConvertToLegacyVote(*resp.Vote)
	if err != nil {
		return nil, err
	}

	return &govv1beta1.QueryVoteResponse{Vote: vote}, nil
}

func (qs legacyQueryServer) Votes(goCtx context.Context, req *govv1beta1.QueryVotesRequest) (*govv1beta1.QueryVotesResponse, error) {
	resp, err := qs.qs.Votes(goCtx, &govv1.QueryVotesRequest{
		ProposalId: req.ProposalId,
		Pagination: req.Pagination,
	})
	if err != nil {
		return nil, err
	}

	votes := make([]govv1beta1.Vote, len(resp.Votes))
	for i, v := range resp.Votes {
		votes[i], err = govv046.ConvertToLegacyVote(*v)
		if err != nil {
			return nil, err
		}
	}

	return &govv1beta1.QueryVotesResponse{
		Votes:      votes,
		Pagination: resp.Pagination,
	}, nil
}

func (qs legacyQueryServer) Params(goCtx context.Context, req *govv1beta1.QueryParamsRequest) (*govv1beta1.QueryParamsResponse, error) {
	resp, err := qs.qs.Params(goCtx, &govv1.QueryParamsRequest{
		ParamsType: req.ParamsType,
	})
	if err != nil {
		return nil, err
	}

	response := &govv1beta1.QueryParamsResponse{}

	if resp.DepositParams != nil {
		minDeposit := sdk.NewCoins(resp.DepositParams.MinDeposit...)
		response.DepositParams = govv1beta1.NewDepositParams(minDeposit, *resp.DepositParams.MaxDepositPeriod)
	}

	if resp.VotingParams != nil {
		response.VotingParams = govv1beta1.NewVotingParams(*resp.VotingParams.VotingPeriod)
	}

	if resp.TallyParams != nil {
		quorum, err := sdk.NewDecFromStr(resp.TallyParams.Quorum)
		if err != nil {
			return nil, err
		}
		threshold, err := sdk.NewDecFromStr(resp.TallyParams.Threshold)
		if err != nil {
			return nil, err
		}
		vetoThreshold, err := sdk.NewDecFromStr(resp.TallyParams.VetoThreshold)
		if err != nil {
			return nil, err
		}

		response.TallyParams = govv1beta1.NewTallyParams(quorum, threshold, vetoThreshold)
	}

	return response, nil
}

func (qs legacyQueryServer) Deposit(goCtx context.Context, req *govv1beta1.QueryDepositRequest) (*govv1beta1.QueryDepositResponse, error) {
	resp, err := qs.qs.Deposit(goCtx, &govv1.QueryDepositRequest{
		ProposalId: req.ProposalId,
		Depositor:  req.Depositor,
	})
	if err != nil {
		return nil, err
	}

	deposit := govv046.ConvertToLegacyDeposit(resp.Deposit)
	return &govv1beta1.QueryDepositResponse{Deposit: deposit}, nil
}

func (qs legacyQueryServer) Deposits(goCtx context.Context, req *govv1beta1.QueryDepositsRequest) (*govv1beta1.QueryDepositsResponse, error) {
	resp, err := qs.qs.Deposits(goCtx, &govv1.QueryDepositsRequest{
		ProposalId: req.ProposalId,
		Pagination: req.Pagination,
	})
	if err != nil {
		return nil, err
	}

	deposits := make([]govv1beta1.Deposit, len(resp.Deposits))
	for idx, deposit := range resp.Deposits {
		deposits[idx] = govv046.ConvertToLegacyDeposit(deposit)
	}

	return &govv1beta1.QueryDepositsResponse{Deposits: deposits, Pagination: resp.Pagination}, nil
}
