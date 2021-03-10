package testutil

import (
	"time"

	"github.com/rafael-sousa/stn-accounts/pkg/model/dto"
	"github.com/rafael-sousa/stn-accounts/pkg/model/entity"
	"github.com/rafael-sousa/stn-accounts/pkg/model/types"
)

// NewEntityAccount returns a new *entity.Account value pointer from the given args
func NewEntityAccount(id int64, n, cpf, s string, b float64) *entity.Account {
	return &entity.Account{
		ID:        id,
		Name:      n,
		CPF:       cpf,
		Secret:    s,
		Balance:   types.NewCurrency(b),
		CreatedAt: time.Now().UTC().Truncate(time.Second),
	}
}

// NewAccountCreation returns a new *dto.AccountCreation value pointer from the given args
func NewAccountCreation(n, cpf, s string, b float64) *dto.AccountCreation {
	return &dto.AccountCreation{
		Name:    n,
		CPF:     cpf,
		Secret:  s,
		Balance: b,
	}
}

// NewEntityTransfer returns a new *entity.Transfer value pointer from the given args
func NewEntityTransfer(id, origin, destination int64, b float64) *entity.Transfer {
	return &entity.Transfer{
		ID:          id,
		Origin:      origin,
		Destination: destination,
		Amount:      types.NewCurrency(b),
		CreatedAt:   time.Now().UTC().Truncate(time.Second),
	}
}
