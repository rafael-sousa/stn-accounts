package testutil

import (
	"time"

	"github.com/rafael-sousa/stn-accounts/pkg/model/entity"
	"github.com/rafael-sousa/stn-accounts/pkg/model/types"
)

// NewAccount returns a new *entity.Account value pointer from the given args
func NewAccount(id int64, n, cpf, s string, b float64) *entity.Account {
	return &entity.Account{
		ID:        id,
		Name:      n,
		CPF:       cpf,
		Secret:    s,
		Balance:   types.NewCurrency(b),
		CreatedAt: time.Now().UTC().Truncate(time.Second),
	}
}
