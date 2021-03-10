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

// NewAccountView returns a new *dto.AccountView value pointer from the given args
func NewAccountView(id int64, name, cpf string, balance float64, createdAt time.Time) *dto.AccountView {
	return &dto.AccountView{
		ID:        id,
		Name:      name,
		CPF:       cpf,
		Balance:   balance,
		CreatedAt: createdAt,
	}
}

// NewTransferCreation returns a new *dto.TransferCreation value pointer from the given args
func NewTransferCreation(dest int64, amt float64) *dto.TransferCreation {
	return &dto.TransferCreation{
		Destination: dest,
		Amount:      amt,
	}
}

// NewTransferView returns a new *dto.TransferView value pointer from the given args
func NewTransferView(id, destination int64, amount float64) *dto.TransferView {
	return &dto.TransferView{
		ID:          id,
		Destination: destination,
		Amount:      amount,
		CreatedAt:   time.Now(),
	}
}
