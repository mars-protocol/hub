package types

import (
	"encoding/json"
)

// ProposalMetadata defines the required schema for proposal metadata.
type ProposalMetadata struct {
	Title             string                 `json:"title"`
	Authors           []string               `json:"authors"`
	Summary           string                 `json:"summary,omitempty"`
	Details           string                 `json:"details"`
	ProposalForumURL  string                 `json:"proposal_forum_url,omitempty"`
	VoteOptionContext string                 `json:"vote_option_context,omitempty"`
	X                 map[string]interface{} `json:"-"` // unexpected fields go here
}

// VoteMetadata defines the required schema for vote metadata.
type VoteMetadata struct {
	Justification string                 `json:"justification,omitempty"`
	X             map[string]interface{} `json:"-"` // unexpected fields go here
}

// UnmarshalProposalMetadata unmarshals a string into ProposalMetadata.
//
// Golang's JSON unmarshal function is retarded. It doesn't check for missing or
// redundant fields. For example, here in ProposalMetadata the "title" field is
// required. But if we provide a JSON string that doesn't have a title, the
// json.Unmarshal simply returns metadata.Title = "". Similarly if the JSON
// string contains an unexpected field, it doesn't throw an error.
//
// Therefore we have to implement our own unmarshal function. See:
//   - Assert required fields are included:
//     https://stackoverflow.com/questions/19633763/unmarshaling-json-in-go-required-field
//   - Assert unknown fields are not included:
//     https://stackoverflow.com/questions/33436730/unmarshal-json-with-some-known-and-some-unknown-field-names
func UnmarshalProposalMetadata(metadataStr string) (*ProposalMetadata, error) {
	var metadata ProposalMetadata

	if err := json.Unmarshal([]byte(metadataStr), &metadata); err != nil {
		return nil, ErrInvalidMetadata.Wrap(err.Error())
	}

	if err := json.Unmarshal([]byte(metadataStr), &metadata.X); err != nil {
		return nil, ErrInvalidMetadata.Wrap(err.Error())
	}

	delete(metadata.X, "title")
	delete(metadata.X, "authors")
	delete(metadata.X, "summary")
	delete(metadata.X, "details")
	delete(metadata.X, "proposal_forum_url")
	delete(metadata.X, "vote_option_context")

	if metadata.Title == "" {
		return nil, ErrInvalidMetadata.Wrap("missing field `title`")
	}

	if metadata.Authors == nil {
		return nil, ErrInvalidMetadata.Wrap("missing field `authors`")
	}

	if metadata.Details == "" {
		return nil, ErrInvalidMetadata.Wrap("missing field `details`")
	}

	if len(metadata.X) > 0 {
		return nil, ErrInvalidMetadata.Wrap("unexpected field(s)")
	}

	return &metadata, nil
}

// UnmarshalVoteMetadata unmarshals a string into VoteMetdata.
//
// See the comments for UnmarshalProposalMetadata on why we need to define this
// function instead of using Go's native json.Unmarshal function.
func UnmarshalVoteMetadata(metadataStr string) (*VoteMetadata, error) {
	var metadata VoteMetadata

	if err := json.Unmarshal([]byte(metadataStr), &metadata); err != nil {
		return nil, ErrInvalidMetadata.Wrap(err.Error())
	}

	if err := json.Unmarshal([]byte(metadataStr), &metadata.X); err != nil {
		return nil, ErrInvalidMetadata.Wrap(err.Error())
	}

	delete(metadata.X, "justification")

	if len(metadata.X) > 0 {
		return nil, ErrInvalidMetadata.Wrap("unexpected field(s)")
	}

	return &metadata, nil
}
