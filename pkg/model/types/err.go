package types

import (
	"fmt"

	"github.com/rs/zerolog"
)

// ErrCode distinguishes the errors expected by the application
type ErrCode string

// List of common application errors
const (
	InternalErr      ErrCode = "0000" // InternalErr represents an unpredicted error such as db connection or network failure
	EmptyResultErr   ErrCode = "0010" // EmptyResultErr represents a state where a given query returned no rows
	SelectStmtErr    ErrCode = "0020" // SelectStmtErr represents a state where a select stmt completed unsuccessfully
	InsertStmtErr    ErrCode = "0030" // InsertStmtErr represents a state where a insert stmt completed unsuccessfully
	UpdateStmtErr    ErrCode = "0040" // UpdateStmtErr represents a state where a update stmt completed unsuccessfully
	NoRowAffectedErr ErrCode = "0050" // NoRowAffectedErr represents a state where a update haven't affect a single row
)

// Err represents an error acknowledged by the application business
type Err struct {
	Code  ErrCode
	Msg   string
	Cause *error
}

var _ zerolog.LogObjectMarshaler = (*Err)(nil)

// MarshalZerologObject appends the current error values to zerolog event logger
func (e *Err) MarshalZerologObject(evt *zerolog.Event) {
	evt.Str("code", string(e.Code)).Str("msg", e.Msg)
	if e.Cause != nil {
		evt.Err(*e.Cause)
	}

}

func (e *Err) Error() string {
	if e.Cause == nil {
		return fmt.Sprintf("%s: %s", e.Code, e.Msg)
	}
	return fmt.Sprintf("%s: %s, %v", e.Code, e.Msg, *e.Cause)
}

// NewErr returns a Err types with the given parameters
func NewErr(c ErrCode, msg string, err *error) *Err {
	return &Err{
		Code:  c,
		Msg:   msg,
		Cause: err,
	}
}
