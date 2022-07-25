package utils

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"

	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govclientrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
)

// GetProposalRESTHandler creates a REST handler for the governance proposal of the specified subRoute
func GetProposalRESTHandler(subRoute string) govclient.RESTHandlerFn {
	return func(client.Context) govclientrest.ProposalRESTHandler {
		return govclientrest.ProposalRESTHandler{
			SubRoute: subRoute,
			Handler:  func(w http.ResponseWriter, r *http.Request) {}, // deprecated, do nothing
		}
	}
}
