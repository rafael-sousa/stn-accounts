package types

import (
	"fmt"

	"github.com/rs/zerolog"
)

// ErrCode distinguishes the errors expected by the application
type ErrCode string

// List of common application errors
const (
	InternalErr       ErrCode = "0000" // InternalErr represents an unpredicted error such as db connection or network failure
	EmptyResultErr    ErrCode = "0010" // EmptyResultErr represents a state where a given query returned no rows
	SelectStmtErr     ErrCode = "0020" // SelectStmtErr represents a state where a select stmt completed unsuccessfully
	InsertStmtErr     ErrCode = "0030" // InsertStmtErr represents a state where a insert stmt completed unsuccessfully
	UpdateStmtErr     ErrCode = "0040" // UpdateStmtErr represents a state where a update stmt completed unsuccessfully
	NoRowAffectedErr  ErrCode = "0050" // NoRowAffectedErr represents a state where a update haven't affect a single row
	ValidationErr     ErrCode = "0060" // ValidationErr represents a business rule violation
	NotFoundErr       ErrCode = "0070" // NotFoundErr occurs when accessing a nonexisting resource
	AuthenticationErr ErrCode = "0080" // AuthenticationErr occurs when the authentication process completes unsuccessfully
	ConflictErr       ErrCode = "0090" // ConflictErr occurs an operation could not complete due to a conflict with the current state of the resource
)

// Err represents an error acknowledged by the application business
type Err struct {
	Code  ErrCode
	Msg   string
	Cause *error
}

var _ zerolog.LogObjectMarshaler = (*Err)(nil)

// NewErr returns a Err types with the given parameters
func NewErr(c ErrCode, msg string, err error) error {
	return &Err{
		Code:  c,
		Msg:   msg,
		Cause: &err,
	}
}

// MarshalZerologObject appends the current error values to zerolog event logger
func (e *Err) MarshalZerologObject(evt *zerolog.Event) {
	evt.Str("code", string(e.Code)).Str("msg", e.Msg)
	if e.Cause != nil {
		evt.Err(*e.Cause)
	}

}

// Error formats a string that describes the custom error
func (e *Err) Error() string {
	if e.Cause == nil {
		return fmt.Sprintf("%s: %s", e.Code, e.Msg)
	}
	return fmt.Sprintf("%s: %s, %v", e.Code, e.Msg, *e.Cause)
}
