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
	AccountRepo *repository.Account
}

// Creation validates the creation of a new entity.Account
func (v *Account) Creation(ctx context.Context, d *dto.AccountCreation) (err error) {
	if err = verifyName(&d.Name); err != nil {
		return
	}
	if err = verifyCPF(&d.CPF); err != nil {
		return
	}
	if err = verifySecret(&d.Secret); err != nil {
		return
	}

	if acc, err := (*v.AccountRepo).FindBy(ctx, d.CPF); acc != nil && err == nil {
		return uniqErr("cpf", d.CPF)
	}
	if d.Balance < 0 {
		return greaterOrEqualErr("balance", 0)
	}
	return
}

// Login validates the creation of a new entity.Account
func (v *Account) Login(cpf *string, secret *string) (err error) {
	if err = verifyCPF(cpf); err != nil {
		return
	}
	nLen := len(*secret)
	if nLen == 0 {
		return requiredFieldErr("secret")
	}
	return
}

func verifyName(v *string) error {
	fieldName := "name"
	nLen := len(*v)
	if nLen == 0 {
		return requiredFieldErr(fieldName)
	}
	if nLen > entity.AccountNameSize {
		return maxSizeErr(fieldName, entity.AccountNameSize)
	}
	if len(strings.TrimSpace(*v)) != nLen {
		return trailingWhiteSpaceErr(fieldName)
	}
	return nil
}

func verifyCPF(v *string) error {
	fieldName := "cpf"
	nLen := len(*v)
	if nLen == 0 {
		return requiredFieldErr(fieldName)
	}
	if nLen > entity.AccountCPFSize {
		return maxSizeErr(fieldName, entity.AccountCPFSize)
	}
	if nLen < entity.AccountCPFSize {
		return minSizeErr(fieldName, entity.AccountCPFSize)
	}
	if _, err := strconv.ParseUint(*v, 10, 64); err != nil {
		return invalidFormatErr(fieldName)
	}
	return nil
}

func verifySecret(v *string) error {
	fieldName := "secret"
	nLen := len(*v)
	if nLen == 0 {
		return requiredFieldErr(fieldName)
	}
	if nLen > entity.AccountSecretSize {
		return maxSizeErr(fieldName, entity.AccountSecretSize)
	}
	return nil
}
