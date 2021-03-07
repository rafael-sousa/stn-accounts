// Package body includes models that handles request and response data
package body

import "time"

// JSONError contains the default error template as json
type JSONError struct {
	Name    string    `json:"name"`
	Code    string    `json:"code"`
	Message string    `json:"message"`
	Path    string    `json:"path"`
	Time    time.Time `json:"time"`
}
