package testutil

import (
	"testing"

	"github.com/rafael-sousa/stn-accounts/pkg/model/types"
)

// AssertCustomErr asserts that the err has the code and msg expected
func AssertCustomErr(t *testing.T, c types.ErrCode, err error, msg string) {
	if customErr, ok := err.(*types.Err); ok {
		AssertEq(t, "err code", c, customErr.Code)
		AssertEq(t, "err msg", msg, customErr.Msg)
	} else {
		t.Errorf("expected err equal to a types.Err but got %v", err)
	}
}

// AssertEq asserts that two values are equal
func AssertEq(t *testing.T, n string, expected interface{}, current interface{}) {
	if expected != current {
		t.Errorf("expected %s equal to '%v' but got '%v'", n, expected, current)
	}
}
