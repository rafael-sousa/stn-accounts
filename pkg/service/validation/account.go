package validation

import (
	"context"
	"strconv"
	"strings"

	"github.com/rafael-sousa/stn-accounts/pkg/model/dto"
	"github.com/rafael-sousa/stn-accounts/pkg/model/entity"
	"github.com/rafael-sousa/stn-accounts/pkg/repository"
)

// Account keeps the validation for operations related to entity.Account
type Account struct {
	AccountRepository *repository.Account
}

// Creation validates the creation of a new entity.Account
func (v *Account) Creation(ctx context.Context, accountCreation *dto.AccountCreation) error {
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

	if account, err := (*v.AccountRepository).FindBy(ctx, accountCreation.CPF); account != nil && err == nil {
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
	cpfLength := len(cpf)
	_, err := strconv.ParseUint(cpf, 10, 64)
	switch {
	case cpfLength == 0:
		return requiredFieldErr(fieldName)
	case cpfLength > entity.AccountCPFSize:
		return maxSizeErr(fieldName, entity.AccountCPFSize)
	case cpfLength < entity.AccountCPFSize:
		return minSizeErr(fieldName, entity.AccountCPFSize)
	case err != nil:
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
