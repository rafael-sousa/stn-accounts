package response

import "net/http"

// List of common header values used by the application
const (
	ApplicationJSON string = "application/json"
	CharsetUTF8     string = "utf-8"
)

// AppendHeaders appends the rest default headers
func AppendHeaders(hs http.Header) {
	hs.Set("Content-Type", ApplicationJSON)
	hs.Set("Accept", ApplicationJSON)
	hs.Set("Accept-Charset", CharsetUTF8)

}
