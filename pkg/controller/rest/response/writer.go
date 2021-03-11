// Package response holds utility functions related to http response handling
package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rafael-sousa/stn-accounts/pkg/controller/rest/body"
	"github.com/rafael-sousa/stn-accounts/pkg/model/types"
)

func appendHeaders(header http.Header) {
	header.Set("Content-Type", "application/json")
	header.Set("Accept", "application/json")
	header.Set("Accept-Charset", "utf-8")
}

// WriteErr writes an appropriate error message into the http response.
// It'll write a respose status code from either 4xx or 5xx category
func WriteErr(w http.ResponseWriter, r *http.Request, err error) error {
	status, errCode := http.StatusInternalServerError, types.InternalErr
	if err, ok := err.(*types.Err); ok {
		errCode = err.Code
		switch err.Code {
		case types.NotFoundErr:
			status = http.StatusNotFound
		case types.EmptyResultErr:
			status = http.StatusNotFound
		case types.ValidationErr:
			status = http.StatusBadRequest
		case types.AuthenticationErr:
			status = http.StatusUnauthorized
		case types.ConflictErr:
			status = http.StatusConflict
		}
	}

	appendHeaders(w.Header())
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

// WriteSuccess writes to the response with appropriate status and headers
func WriteSuccess(w http.ResponseWriter, r *http.Request, b interface{}, id interface{}) error {
	appendHeaders(w.Header())
	if b == nil {
		w.WriteHeader(http.StatusNoContent)
		return nil
	}
	status := http.StatusOK
	if id != nil {
		status = http.StatusCreated
		w.Header().Set("Location", fmt.Sprintf("%s/%v", r.URL.RequestURI(), id))
	}
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(b)
}
