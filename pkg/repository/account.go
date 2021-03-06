// Package repository exposes interfaces meant to serve outer layers and general artefacts there aren't implementation-specific.
// Its subpackages contain implementations to specific databases
package repository

import (
	"context"

	"github.com/rafael-sousa/stn-accounts/pkg/model/entity"
	"github.com/rafael-sousa/stn-accounts/pkg/model/types"
)

// Account exposes database operations related to account domain
type Account interface {
	Fetch(ctx context.Context) ([]*entity.Account, error)
	Create(ctx context.Context, e *entity.Account) (*entity.Account, error)
	GetBalance(ctx context.Context, id int64) (types.Currency, error)
	FindBy(ctx context.Context, cpf string) (*entity.Account, error)
	UpdateBalance(ctx context.Context, id int64, b types.Currency) error
}
