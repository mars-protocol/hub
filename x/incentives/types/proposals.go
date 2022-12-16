package types

import (
	"fmt"

	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	marsutils "github.com/mars-protocol/hub/utils"
)

const (
	ProposalCreateIncentivesSchedule   = "IncentivesScheduleCreation"
	ProposalTerminalIncentivesSchedule = "IncentivesScheduleTermination"
)

var (
	_ govv1beta1.Content = &CreateIncentivesScheduleProposal{}
	_ govv1beta1.Content = &TerminateIncentivesSchedulesProposal{}
)

func init() {
	govv1beta1.RegisterProposalType(ProposalCreateIncentivesSchedule)
	govv1beta1.RegisterProposalType(ProposalTerminalIncentivesSchedule)
}

//--------------------------------------------------------------------------------------------------
// CreateIncentivesScheduleProposal
//--------------------------------------------------------------------------------------------------

func (p *CreateIncentivesScheduleProposal) GetTitle() string {
	return p.Title
}

func (p *CreateIncentivesScheduleProposal) GetDescription() string {
	return p.Description
}

func (p *CreateIncentivesScheduleProposal) ProposalRoute() string {
	return RouterKey
}

func (p *CreateIncentivesScheduleProposal) ProposalType() string {
	return ProposalCreateIncentivesSchedule
}

func (p *CreateIncentivesScheduleProposal) ValidateBasic() error {
	if err := govv1beta1.ValidateAbstract(p); err != nil {
		return err
	}

	if !p.StartTime.Before(p.EndTime) {
		return ErrInvalidProposalStartEndTimes
	}

	if p.Amount.Empty() {
		return ErrInvalidProposalAmount
	}

	return nil
}

func (p CreateIncentivesScheduleProposal) String() string {
	return fmt.Sprintf(
		`Incentives Schedule Creation Proposal:
  Title:       %s
  Description: %s
  Start time:  %s
  End time:    %s
  Amount:      %s
`,
		p.Title,
		p.Description,
		p.StartTime,
		p.EndTime,
		p.Amount,
	)
}

//--------------------------------------------------------------------------------------------------
// TerminateIncentivesSchedulesProposal
//--------------------------------------------------------------------------------------------------

func (p *TerminateIncentivesSchedulesProposal) GetTitle() string {
	return p.Title
}

func (p *TerminateIncentivesSchedulesProposal) GetDescription() string {
	return p.Description
}

func (p *TerminateIncentivesSchedulesProposal) ProposalRoute() string {
	return RouterKey
}

func (p *TerminateIncentivesSchedulesProposal) ProposalType() string {
	return ProposalTerminalIncentivesSchedule
}

func (p *TerminateIncentivesSchedulesProposal) ValidateBasic() error {
	if err := govv1beta1.ValidateAbstract(p); err != nil {
		return err
	}

	if len(p.Ids) == 0 {
		return ErrInvalidProposalIds
	}

	return nil
}

func (p TerminateIncentivesSchedulesProposal) String() string {
	return fmt.Sprintf(
		`Incentives Schedule Termination Proposal:
  Title:       %s
  Description: %s
  Ids:         %s
`,
		p.Title,
		p.Description,
		marsutils.UintArrayToString(p.Ids, ", "),
	)
}
