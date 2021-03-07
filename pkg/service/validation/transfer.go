package validation

import (
	"context"
	"fmt"

	"github.com/rafael-sousa/stn-accounts/pkg/model/dto"
	"github.com/rafael-sousa/stn-accounts/pkg/model/types"
	"github.com/rafael-sousa/stn-accounts/pkg/repository"
)

// Transfer keeps the validation for operations related to entity.Transfer
type Transfer struct {
	AccountRepo *repository.Account
}

// Creation validates the creation of a new entity.Transfer
func (v *Transfer) Creation(ctx context.Context, origin int64, d *dto.TransferCreation) error {
	if d.Amount <= 0 {
		return greaterThanErr("amount", 0)
	}
	if d.Destination == 0 {
		return requiredFieldErr("destination_id")
	}
	if d.Destination == origin {
		return sameFieldErr("origin id", "destination id")
	}
	oBalance, err := (*v.AccountRepo).GetBalance(ctx, origin)
	if err != nil {
		return notFoundErr("origin", origin)
	}
	b := int64(oBalance) - int64(d.Amount*100)
	if b < 0 {
		return types.NewErr(types.ValidationErr, fmt.Sprintf("the origin must have a balance greater than or equal to %.2f", d.Amount), nil)
	}
	_, err = (*v.AccountRepo).GetBalance(ctx, d.Destination)
	if err != nil {
		return notFoundErr("destination", origin)
	}
	return nil
}
