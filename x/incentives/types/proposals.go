package types

import (
	"fmt"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	marsutils "github.com/mars-protocol/hub/utils"
)

const (
	ProposalCreateIncentivesSchedule   = "IncentivesScheduleCreation"
	ProposalTerminalIncentivesSchedule = "IncentivesScheduleTermination"
)

var (
	_ govtypes.Content = &CreateIncentivesScheduleProposal{}
	_ govtypes.Content = &TerminateIncentivesSchedulesProposal{}
)

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
	if err := govtypes.ValidateAbstract(p); err != nil {
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
	if err := govtypes.ValidateAbstract(p); err != nil {
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
