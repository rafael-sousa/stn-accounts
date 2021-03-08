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
	Fetch(ctx context.Context, id int64) ([]*dto.TransferView, error)
	Create(ctx context.Context, origin int64, d *dto.TransferCreation) (*dto.TransferView, error)
}

type transfer struct {
	transferRepo *repository.Transfer
	accountRepo  *repository.Account
	txr          *repository.Transactioner
	transferVali *validation.Transfer
}

// Fetch returns a list of entity.Transfer from the entity.Account stored at id.
// It returns nil and an error in when not able to fetch the rows from the repository
func (s *transfer) Fetch(ctx context.Context, id int64) ([]*dto.TransferView, error) {
	transfers, err := (*s.transferRepo).Fetch(ctx, id)
	if err != nil {
		log.Error().Caller().Err(err).Int64("id", id).Msg("unable to fetch transfers")
		return nil, err
	}

	views := make([]*dto.TransferView, 0, len(transfers))
	for _, t := range transfers {
		views = append(views, dto.NewTransferView(t))
	}
	return views, nil
}

// Create validates, create, and persists an entity.Transfer from the values stored at d
func (s *transfer) Create(ctx context.Context, origin int64, d *dto.TransferCreation) (*dto.TransferView, error) {
	var e *entity.Transfer
	err := (*s.txr).WithTx(ctx, func(txCtx context.Context) error {
		err := s.transferVali.Creation(txCtx, origin, d)
		if err != nil {
			return err
		}
		originBalance, err := (*s.accountRepo).GetBalance(txCtx, origin)
		if err != nil {
			log.Info().Caller().Err(err).Int64("account_origin_id", origin).Msg("unable to get the origin account balance")
			return err
		}
		destBalance, err := (*s.accountRepo).GetBalance(txCtx, d.Destination)
		if err != nil {
			log.Info().Caller().Err(err).Int64("account_destination_id", d.Destination).Msg("unable to get the destination account balance")
			return err
		}
		amount := types.NewCurrency(d.Amount)
		err = (*s.accountRepo).UpdateBalance(txCtx, origin, originBalance-amount)
		if err != nil {
			log.Info().Caller().Err(err).
				Int64("account_origin_id", origin).
				Int64("balance", int64(originBalance)).
				Int64("amount", int64(amount)).
				Msg("unable to update the origin account balance")
			return err
		}
		err = (*s.accountRepo).UpdateBalance(txCtx, d.Destination, destBalance+amount)
		if err != nil {
			log.Info().
				Caller().
				Err(err).
				Int64("account_destination_id", d.Destination).
				Int64("balance", int64(destBalance)).
				Int64("amount", int64(amount)).
				Msg("unable to update the destination account balance")
			return err
		}
		e, err = (*s.transferRepo).Create(txCtx, &entity.Transfer{
			Origin:      origin,
			Destination: d.Destination,
			Amount:      amount,
			CreatedAt:   time.Now(),
		})

		return err
	})
	if err != nil {
		log.Info().
			Caller().
			Err(err).
			Int64("account_origin_id", origin).
			Int64("account_destination_id", d.Destination).
			Int64("amount", int64(d.Amount)).
			Msg("unable to transfer the currency amount")
		return nil, err
	}
	return dto.NewTransferView(e), nil
}

// NewTransfer returns a value responsible for managing entity.Transfer integrity
func NewTransfer(txr *repository.Transactioner, transferRepo *repository.Transfer, accountRepo *repository.Account) Transfer {
	return &transfer{
		transferRepo: transferRepo,
		accountRepo:  accountRepo,
		txr:          txr,
		transferVali: &validation.Transfer{
			AccountRepo: accountRepo,
		},
	}
}
