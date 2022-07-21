package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const ProposalTypeSafetyFundSpend = "SafetyFundSpend"

var _ govtypes.Content = &SafetyFundSpendProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypeSafetyFundSpend)
	govtypes.RegisterProposalTypeCodec(&SafetyFundSpendProposal{}, "mars/SafetyFundSpendProposal")
}

// NewSafetyFundSpendProposal creates a new instance of SafetyFundSpendProposal
func NewSafetyFundSpendProposal(title, description string, recipientAddr sdk.AccAddress, amount sdk.Coins) *SafetyFundSpendProposal {
	return &SafetyFundSpendProposal{
		Title:       title,
		Description: description,
		Recipient:   recipientAddr.String(),
		Amount:      amount,
	}
}

func (sfsp *SafetyFundSpendProposal) GetTitle() string {
	return sfsp.Title
}

func (sfsp *SafetyFundSpendProposal) GetDescription() string {
	return sfsp.Description
}

func (sfsp *SafetyFundSpendProposal) ProposalRoute() string {
	return RouterKey
}

func (sfsp *SafetyFundSpendProposal) ProposalType() string {
	return ProposalTypeSafetyFundSpend
}

func (sfsp *SafetyFundSpendProposal) ValidateBasic() error {
	if err := govtypes.ValidateAbstract(sfsp); err != nil {
		return err
	}

	if !sfsp.Amount.IsValid() {
		return ErrInvalidProposalAmount
	}

	if sfsp.Recipient == "" {
		return ErrEmptyProposalRecipient
	}

	return nil
}

func (sfsp SafetyFundSpendProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(
		`Community Pool Spend Proposal:
  Title:       %s
  Description: %s
  Recipient:   %s
  Amount:      %s
`,
		sfsp.Title,
		sfsp.Description,
		sfsp.Recipient,
		sfsp.Amount,
	))
	return b.String()
}
