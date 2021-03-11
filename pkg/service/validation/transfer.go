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
	AccountRepository *repository.Account
}

// Creation validates the creation of a new entity.Transfer
func (v *Transfer) Creation(ctx context.Context, origin int64, transferCreation dto.TransferCreation) error {
	if transferCreation.Amount <= 0 {
		return greaterThanErr("amount", 0)
	}
	if transferCreation.Destination <= 0 {
		return requiredFieldErr("destination_id")
	}
	if transferCreation.Destination == origin {
		return sameFieldErr("origin id", "destination id")
	}
	originBalance, err := (*v.AccountRepository).GetBalance(ctx, origin)
	if err != nil {
		if customErr, ok := err.(*types.Err); ok && customErr.Code == types.EmptyResultErr {
			return notFoundErr("origin", origin)
		}
		return err
	}
	if originBalance-types.NewCurrency(transferCreation.Amount) < 0 {
		return types.NewErr(types.ValidationErr, fmt.Sprintf("the origin must have a balance greater than or equal to %.2f", transferCreation.Amount), nil)
	}
	exists, err := (*v.AccountRepository).Exists(ctx, transferCreation.Destination)
	if err != nil {
		return err
	} else if !exists {
		return notFoundErr("destination", transferCreation.Destination)
	}
	return nil
}
