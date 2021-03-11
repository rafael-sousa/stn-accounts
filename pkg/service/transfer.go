package service

import (
	"context"
	"time"

	"github.com/rafael-sousa/stn-accounts/pkg/model/dto"
	"github.com/rafael-sousa/stn-accounts/pkg/model/entity"
	"github.com/rafael-sousa/stn-accounts/pkg/model/types"
	"github.com/rafael-sousa/stn-accounts/pkg/repository"
	"github.com/rafael-sousa/stn-accounts/pkg/service/validation"
	"github.com/rs/zerolog/log"
)

// Transfer represents the business operations available to entity.Transfer type
type Transfer interface {
	Fetch(ctx context.Context, id int64) ([]dto.TransferView, error)
	Create(ctx context.Context, origin int64, d dto.TransferCreation) (dto.TransferView, error)
}

type transfer struct {
	transferRepository *repository.Transfer
	accountRepository  *repository.Account
	txr                *repository.Transactioner
	transferValidator  *validation.Transfer
}

// Fetch returns a list of entity.Transfer from the entity.Account stored at id.
// It returns nil and an error in when not able to fetch the rows from the repository
func (s *transfer) Fetch(ctx context.Context, id int64) ([]dto.TransferView, error) {
	transfers, err := (*s.transferRepository).Fetch(ctx, id)
	if err != nil {
		log.Error().Caller().Err(err).Int64("id", id).Msg("unable to fetch transfers")
		return nil, err
	}

	views := make([]dto.TransferView, 0, len(transfers))
	for _, t := range transfers {
		views = append(views, dto.NewTransferView(t))
	}
	return views, nil
}

// Create validates, create, and persists an entity.Transfer from the values stored at d
func (s *transfer) Create(ctx context.Context, origin int64, transferCreation dto.TransferCreation) (view dto.TransferView, err error) {
	var transfer entity.Transfer
	err = (*s.txr).WithTx(ctx, func(txCtx context.Context) error {
		if err := s.transferValidator.Creation(txCtx, origin, transferCreation); err != nil {
			return err
		}
		originBalance, err := (*s.accountRepository).GetBalance(txCtx, origin)
		if err != nil {
			log.Info().Caller().Err(err).
				Int64("account_origin_id", origin).
				Msg("unable to get the origin account balance")
			return err
		}
		destinationBalance, err := (*s.accountRepository).GetBalance(txCtx, transferCreation.Destination)
		if err != nil {
			log.Info().Caller().Err(err).
				Int64("account_destination_id", transferCreation.Destination).
				Msg("unable to get the destination account balance")
			return err
		}
		amount := types.NewCurrency(transferCreation.Amount)

		if err = (*s.accountRepository).UpdateBalance(txCtx, origin, originBalance-amount); err != nil {
			log.Info().Caller().Err(err).
				Int64("account_origin_id", origin).
				Int64("balance", int64(originBalance)).
				Int64("amount", int64(amount)).
				Msg("unable to update the origin account balance")
			return err
		}

		if err = (*s.accountRepository).UpdateBalance(txCtx, transferCreation.Destination, destinationBalance+amount); err != nil {
			log.Info().
				Caller().
				Err(err).
				Int64("account_destination_id", transferCreation.Destination).
				Int64("balance", int64(destinationBalance)).
				Int64("amount", int64(amount)).
				Msg("unable to update the destination account balance")
			return err
		}
		transfer = entity.Transfer{
			Origin:      origin,
			Destination: transferCreation.Destination,
			Amount:      amount,
			CreatedAt:   time.Now(),
		}
		id, err := (*s.transferRepository).Create(txCtx, transfer)
		transfer.ID = id
		return err
	})
	if err != nil {
		log.Info().
			Caller().
			Err(err).
			Int64("account_origin_id", origin).
			Int64("account_destination_id", transferCreation.Destination).
			Int64("amount", int64(transferCreation.Amount)).
			Msg("unable to transfer the currency amount")
		return view, err
	}
	return dto.NewTransferView(transfer), nil
}

// NewTransfer returns a value responsible for managing entity.Transfer integrity
func NewTransfer(txr *repository.Transactioner, transferRepository *repository.Transfer, accountRepository *repository.Account) Transfer {
	return &transfer{
		transferRepository: transferRepository,
		accountRepository:  accountRepository,
		txr:                txr,
		transferValidator: &validation.Transfer{
			AccountRepository: accountRepository,
		},
	}
}
