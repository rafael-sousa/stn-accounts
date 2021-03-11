// Package validation groups extensive business validation rules
package validation

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/rafael-sousa/stn-accounts/pkg/model/dto"
	"github.com/rafael-sousa/stn-accounts/pkg/model/entity"
	"github.com/rafael-sousa/stn-accounts/pkg/repository"
)

var _hasElevenDigits = regexp.MustCompile(`^([\d]{11})$`).MatchString

// Account keeps the validation for operations related to entity.Account
type Account struct {
	AccountRepository *repository.Account
}

// Creation validates the creation of a new entity.Account
func (v *Account) Creation(ctx context.Context, accountCreation dto.AccountCreation) error {
	if err := verifyName(accountCreation.Name); err != nil {
		return err
	}
	if err := verifyCPF(accountCreation.CPF); err != nil {
		return err
	}
	if err := verifySecret(accountCreation.Secret); err != nil {
		return err
	}
	if accountCreation.Balance < 0 {
		return greaterOrEqualErr("balance", 0)
	}

	if _, err := (*v.AccountRepository).FindBy(ctx, accountCreation.CPF); err == nil {
		return uniqErr("cpf", accountCreation.CPF)
	}

	return nil
}

// Login validates the creation of a new entity.Account
func (v *Account) Login(cpf string, secret string) error {
	if err := verifyCPF(cpf); err != nil {
		return err
	}
	if len(secret) == 0 {
		return requiredFieldErr("secret")
	}
	return nil
}

func verifyName(name string) error {
	fieldName := "name"
	nameLength := len(name)
	switch {
	case nameLength == 0:
		return requiredFieldErr(fieldName)
	case nameLength > entity.AccountNameSize:
		return maxSizeErr(fieldName, entity.AccountNameSize)
	case len(strings.TrimSpace(name)) != nameLength:
		return trailingWhiteSpaceErr(fieldName)
	}
	return nil
}

func verifyCPF(cpf string) error {
	fieldName := "cpf"
	if len(cpf) == 0 {
		return requiredFieldErr(fieldName)
	}

	if !_hasElevenDigits(cpf) {
		return invalidFormatErr(fieldName)
	}

	runes := make(map[rune]rune, 11)
	for _, c := range cpf {
		runes[c] = c
	}
	if len(runes) == 1 {
		return invalidFormatErr(fieldName)
	}

	var s1, s2 int
	for i, c := range cpf[0:9] {
		d, _ := strconv.Atoi(string(c))
		s1 += (10 - i) * d
		s2 += d
	}
	s2 += s1
	var m1, m2 int
	if s1%11 >= 2 {
		m1 = 11 - s1%11
	}
	s2 += m1 * 2
	if s2%11 >= 2 {
		m2 = 11 - s2%11
	}
	if cpf[9:11] != fmt.Sprintf("%d%d", m1, m2) {
		return invalidFormatErr(fieldName)
	}
	return nil
}

func verifySecret(secret string) error {
	fieldName := "secret"
	nLen := len(secret)
	switch {
	case nLen == 0:
		return requiredFieldErr(fieldName)
	case nLen > entity.AccountSecretSize:
		return maxSizeErr(fieldName, entity.AccountSecretSize)
	}
	return nil
}
