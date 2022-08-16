package types

import (
	"fmt"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalExecuteRemoteContract = "RemoteContractExecution"
	ProposalMigrateRemoteContract = "RemoveContractMigration"
)

var (
	_ govtypes.Content = &ExecuteRemoteContractProposal{}
	_ govtypes.Content = &MigrateRemoteContractProposal{}
)

func init() {
	govtypes.RegisterProposalType(ProposalExecuteRemoteContract)
	govtypes.RegisterProposalTypeCodec(&ExecuteRemoteContractProposal{}, "mars/ExecuteRemoteContractProposal")

	govtypes.RegisterProposalType(ProposalMigrateRemoteContract)
	govtypes.RegisterProposalTypeCodec(&MigrateRemoteContractProposal{}, "mars/MigrateRemoteContractProposal")
}

//--------------------------------------------------------------------------------------------------
// ExecuteRemoteContractProposal
//--------------------------------------------------------------------------------------------------

func (p *ExecuteRemoteContractProposal) GetTitle() string {
	return p.Title
}

func (p *ExecuteRemoteContractProposal) GetDescription() string {
	return p.Description
}

func (p *ExecuteRemoteContractProposal) ProposalRoute() string {
	return RouterKey
}

func (p *ExecuteRemoteContractProposal) ProposalType() string {
	return ProposalExecuteRemoteContract
}

func (p *ExecuteRemoteContractProposal) ValidateBasic() error {
	if err := govtypes.ValidateAbstract(p); err != nil {
		return err
	}

	if err := p.ExecuteMsg.ValidateBasic(); err != nil {
		return err
	}

	return nil
}

func (p *ExecuteRemoteContractProposal) String() string {
	msg, err := p.ExecuteMsg.MarshalJSON()
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf(
		`Execute Remote Contract Proposal:
  Title:       %s
  Description: %s
  Chain ID:    %s
  Contract:    %s
  Message:     %s
`,
		p.Title,
		p.Description,
		p.ChainID,
		p.Contract,
		string(msg),
	)
}

//--------------------------------------------------------------------------------------------------
// MigrateRemoteContractProposal
//--------------------------------------------------------------------------------------------------

func (p *MigrateRemoteContractProposal) GetTitle() string {
	return p.Title
}

func (p *MigrateRemoteContractProposal) GetDescription() string {
	return p.Description
}

func (p *MigrateRemoteContractProposal) ProposalRoute() string {
	return RouterKey
}

func (p *MigrateRemoteContractProposal) ProposalType() string {
	return ProposalMigrateRemoteContract
}

func (p *MigrateRemoteContractProposal) ValidateBasic() error {
	if err := govtypes.ValidateAbstract(p); err != nil {
		return err
	}

	if err := p.MigrateMsg.ValidateBasic(); err != nil {
		return err
	}

	return nil
}

func (p *MigrateRemoteContractProposal) String() string {
	msg, err := p.MigrateMsg.MarshalJSON()
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf(
		`Migrate Remote Contract Proposal:
  Title:       %s
  Description: %s
  Chain ID:    %s
  Contract:    %s
  Message:     %s
`,
		p.Title,
		p.Description,
		p.ChainID,
		p.Contract,
		string(msg),
	)
}
