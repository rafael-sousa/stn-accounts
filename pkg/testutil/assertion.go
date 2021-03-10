package testutil

import (
	"testing"
	"time"

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

// AssertNotDefault asserts that the value store at v does not contains its default type value
func AssertNotDefault(t *testing.T, n string, v interface{}) {
	switch v.(type) {
	case time.Time:
		if v.(time.Time).IsZero() {
			t.Errorf("expected %s not empty", n)
		}
	case string:
		if len(v.(string)) == 0 {
			t.Errorf("expected %s not empty", n)
		}
	case int, int32, int64, float32, float64:
		if v == 0 {
			t.Errorf("expected %s not zero", n)
		}
	case nil:
		t.Errorf("expected %s not nil", n)
	}
}
