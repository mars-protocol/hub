package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

const ProposalTypeSafetyFundSpend = "SafetyFundSpend"

var _ govv1beta1.Content = &SafetyFundSpendProposal{}

func init() {
	govv1beta1.RegisterProposalType(ProposalTypeSafetyFundSpend)
	govv1.RegisterLegacyAminoCodec(&SafetyFundSpendProposal{}, "mars/SafetyFundSpendProposal", nil)
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
	if err := govv1beta1.ValidateAbstract(sfsp); err != nil {
		return err
	}

	if !sfsp.Amount.IsValid() {
		return ErrInvalidProposalAmount
	}

	if sfsp.Recipient == "" {
		return ErrInvalidProposalRecipient
	}

	return nil
}

func (sfsp SafetyFundSpendProposal) String() string {
	return fmt.Sprintf(
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
	)
}
