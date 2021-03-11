package repository

import (
	"context"

	"github.com/rafael-sousa/stn-accounts/pkg/model/entity"
)

// Transfer exposes database operations related to transfer domain
type Transfer interface {
	Fetch(ctx context.Context, origin int64) ([]entity.Transfer, error)
	Create(ctx context.Context, e entity.Transfer) (int64, error)
}
