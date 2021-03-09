package response

import "net/http"

// List of common header values used by the application
const (
	ApplicationJSON string = "application/json"
	CharsetUTF8     string = "utf-8"
)

// AppendHeaders appends the rest default headers
func AppendHeaders(header http.Header) {
	header.Set("Content-Type", ApplicationJSON)
	header.Set("Accept", ApplicationJSON)
	header.Set("Accept-Charset", CharsetUTF8)

}
