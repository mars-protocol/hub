package types

import (
	"encoding/json"
)

// ProposalMetadata defines the required schema for proposal metadata.
type ProposalMetadata struct {
	Title             string   `json:"title"`
	Authors           []string `json:"authors,omitempty"`
	Summary           string   `json:"summary"`
	Details           string   `json:"details,omitempty"`
	ProposalForumURL  string   `json:"proposal_forum_url,omitempty"`
	VoteOptionContext string   `json:"vote_option_context,omitempty"`
}

// VoteMetadata defines the required schema for vote metadata.
type VoteMetadata struct {
	Justification string `json:"justification,omitempty"`
}

// UnmarshalProposalMetadata unmarshals a string into ProposalMetadata.
//
// Golang's JSON unmarshal function doesn't check for missing fields. Instead,
// for example, if the "title" field here in ProposalMetadata is missing, the
// json.Unmarshal simply returns metadata.Title = "" instead of throwing an
// error.
//
// Here's the equivalent Rust code for comparison, which properly throws an
// error is a required field is missing:
// https://play.rust-lang.org/?version=stable&mode=debug&edition=2021&gist=0e2eadad38b7cd212962b1a0e7a6da44
//
// Therefore, we have to implement our own unmarshal function which checks for
// missing fields.
func UnmarshalProposalMetadata(metadataStr string) (*ProposalMetadata, error) {
	var metadata ProposalMetadata

	if err := json.Unmarshal([]byte(metadataStr), &metadata); err != nil {
		return nil, ErrInvalidMetadata.Wrap(err.Error())
	}

	if metadata.Title == "" {
		return nil, ErrInvalidMetadata.Wrap("missing field `title`")
	}

	if metadata.Summary == "" {
		return nil, ErrInvalidMetadata.Wrap("missing field `summary`")
	}

	return &metadata, nil
}

// UnmarshalVoteMetadata unmarshals a string into VoteMetdata.
func UnmarshalVoteMetadata(metadataStr string) (*VoteMetadata, error) {
	var metadata VoteMetadata

	if err := json.Unmarshal([]byte(metadataStr), &metadata); err != nil {
		return nil, ErrInvalidMetadata.Wrap(err.Error())
	}

	return &metadata, nil
}
