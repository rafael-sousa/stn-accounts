package response

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/rafael-sousa/stn-accounts/pkg/controller/rest/body"
	"github.com/rafael-sousa/stn-accounts/pkg/model/types"
)

func parseErrorCode(v interface{}) (int, types.ErrCode) {
	if err, ok := v.(*types.Err); ok {
		switch err.Code {
		case types.NotFoundErr:
			return http.StatusNotFound, types.NotFoundErr
		case types.EmptyResultErr:
			return http.StatusNotFound, types.EmptyResultErr
		case types.ValidationErr:
			return http.StatusBadRequest, types.ValidationErr
		case types.AuthenticationErr:
			return http.StatusUnauthorized, types.AuthenticationErr
		case types.ConflictErr:
			return http.StatusConflict, types.ConflictErr
		}
	}

	return http.StatusInternalServerError, types.InternalErr
}

// WriteErr writes an appropriate error message into the http response.
// It'll write a respose status code from either 4xx or 5xx category
func WriteErr(w http.ResponseWriter, r *http.Request, err error) error {
	status, errCode := parseErrorCode(err)
	AppendHeaders(w.Header())
	w.WriteHeader(status)
	msg := err.Error()
	if appErr, ok := err.(*types.Err); ok {
		msg = appErr.Msg
	}
	return json.NewEncoder(w).Encode(&body.JSONError{
		Name:    http.StatusText(status),
		Code:    string(errCode),
		Time:    time.Now(),
		Message: msg,
		Path:    r.URL.RequestURI(),
	})
}
