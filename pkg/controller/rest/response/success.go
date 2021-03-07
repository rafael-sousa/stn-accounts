package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// WriteSuccess writes to the response with appropriate status and headers
func WriteSuccess(w http.ResponseWriter, r *http.Request, b interface{}, id interface{}) error {
	AppendHeaders(w.Header())
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
