package routing

import (
	"net/http"

	"github.com/rafael-sousa/stn-accounts/pkg/controller/rest/response"
	"github.com/rafael-sousa/stn-accounts/pkg/model/types"
)

// NotFound is a route for dealing with nonexisting resources
func NotFound(w http.ResponseWriter, r *http.Request) {
	response.WriteErr(w, r, types.NewErr(types.NotFoundErr, "Unable to match any Request-URI", nil))
}
