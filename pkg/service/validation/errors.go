package validation

import (
	"fmt"

	"github.com/rafael-sousa/stn-accounts/pkg/model/types"
)

func requiredFieldErr(n string) error {
	return types.NewErr(types.ValidationErr, fmt.Sprintf("field '%s' is required", n), nil)
}

func greaterThanErr(n string, v interface{}) error {
	return types.NewErr(types.ValidationErr, fmt.Sprintf("field '%s' must be greater than %v", n, v), nil)
}

func greaterOrEqualErr(n string, v interface{}) error {
	return types.NewErr(types.ValidationErr, fmt.Sprintf("field '%s' must be greater than or equal to %v", n, v), nil)
}

func uniqErr(n string, v interface{}) error {
	return types.NewErr(types.ConflictErr, fmt.Sprintf("field '%s' with value '%v' is already in use", n, v), nil)
}

func notFoundErr(n string, v interface{}) error {
	return types.NewErr(types.EmptyResultErr, fmt.Sprintf("record with '%s' equals '%v' was not found", n, v), nil)
}

func maxSizeErr(n string, s int) error {
	return types.NewErr(types.ValidationErr, fmt.Sprintf("field '%s' must have at most %d characters", n, s), nil)
}

func trailingWhiteSpaceErr(n string) error {
	return types.NewErr(types.ValidationErr, fmt.Sprintf("field '%s' can't have trailing whitespace", n), nil)
}
func invalidFormatErr(n string) error {
	return types.NewErr(types.ValidationErr, fmt.Sprintf("field '%s' has an invalid format", n), nil)
}
func sameFieldErr(n1 string, n2 string) error {
	return types.NewErr(types.ConflictErr, fmt.Sprintf("fields '%s' and '%s' can't be the same", n1, n2), nil)
}
