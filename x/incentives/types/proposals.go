package types

import (
	"fmt"
	"strings"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalCreateIncentivesSchedule   = "IncentivesScheduleCreation"
	ProposalTerminalIncentivesSchedule = "IncentivesScheduleTermination"
)

var (
	_ govtypes.Content = &CreateIncentivesScheduleProposal{}
	_ govtypes.Content = &TerminateIncentivesScheduleProposal{}
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
// TerminateIncentivesScheduleProposal
//--------------------------------------------------------------------------------------------------

func (p *TerminateIncentivesScheduleProposal) GetTitle() string {
	return p.Title
}

func (p *TerminateIncentivesScheduleProposal) GetDescription() string {
	return p.Description
}

func (p *TerminateIncentivesScheduleProposal) ProposalRoute() string {
	return RouterKey
}

func (p *TerminateIncentivesScheduleProposal) ProposalType() string {
	return ProposalTerminalIncentivesSchedule
}

func (p *TerminateIncentivesScheduleProposal) ValidateBasic() error {
	if err := govtypes.ValidateAbstract(p); err != nil {
		return err
	}

	if len(p.Ids) == 0 {
		return ErrInvalidProposalIds
	}

	return nil
}

func (p TerminateIncentivesScheduleProposal) String() string {
	return fmt.Sprintf(
		`Incentives Schedule Termination Proposal:
  Title:       %s
  Description: %s
  Ids:         %s
`,
		p.Title,
		p.Description,
		arrayToString(p.Ids, ", "),
	)
}

// https://stackoverflow.com/questions/37532255/one-liner-to-transform-int-into-string
func arrayToString(a []uint64, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}
