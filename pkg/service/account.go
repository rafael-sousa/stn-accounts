// Package service holds files responsible for input validation and managing application integrity.
// This package contains the application core business.
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
	"golang.org/x/crypto/bcrypt"
)

// Account exposes the business operations available to entity.Account type
type Account interface {
	Fetch(ctx context.Context) ([]*dto.AccountView, error)
	GetBalance(ctx context.Context, id int64) (float64, error)
	Create(ctx context.Context, d *dto.AccountCreation) (*dto.AccountView, error)
	Login(ctx context.Context, cpf string, secret string) (*dto.AccountView, error)
}

type account struct {
	accountRepo *repository.Account
	accountVali *validation.Account
	txr         *repository.Transactioner
}

// Fetch returns a list of dto.AccountView
func (s *account) Fetch(ctx context.Context) ([]*dto.AccountView, error) {
	accs, err := (*s.accountRepo).Fetch(ctx)
	if err != nil {
		log.Error().Caller().Err(err).Msg("unable to fetch accounts")
		return nil, err
	}
	views := make([]*dto.AccountView, 0, len(accs))
	for _, acc := range accs {
		views = append(views, dto.NewAccountView(acc))
	}
	return views, nil
}

// GetBalance returns the given account balance
func (s *account) GetBalance(ctx context.Context, id int64) (float64, error) {
	b, err := (*s.accountRepo).GetBalance(ctx, id)
	if err != nil {
		log.Info().Caller().Err(err).Int64("id", id).Msg("unable to get account balance")
		return 0, err
	}
	return b.Float64(), nil
}

// Create validates and persists the given e entity.Account
func (s *account) Create(ctx context.Context, d *dto.AccountCreation) (*dto.AccountView, error) {
	var e *entity.Account
	err := (*s.txr).WithTx(ctx, func(txCtx context.Context) error {
		err := s.accountVali.Creation(txCtx, d)
		if err != nil {
			return err
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(d.Secret), bcrypt.DefaultCost)
		if err != nil {
			log.Info().Caller().Err(err).Msg("unable to create the account secret hash")
			return err
		}
		e, err = (*s.accountRepo).Create(txCtx, &entity.Account{
			Name:      d.Name,
			CPF:       d.CPF,
			Balance:   types.NewCurrency(d.Balance),
			CreatedAt: time.Now(),
			Secret:    string(hash),
		})
		return err
	})

	if err != nil {
		log.Info().Caller().Err(err).Str("cpf", d.CPF).Str("name", d.Name).Msg("unable to create account")
		return nil, err
	}

	return dto.NewAccountView(e), nil
}

// Login is responsible for fetching an account, comparing the secret, and returning the respective account view.
// It returns nil and an err if the account is not found or the provided secret doesn't match the account
func (s *account) Login(ctx context.Context, cpf string, secret string) (*dto.AccountView, error) {
	err := s.accountVali.Login(&cpf, &secret)
	if err != nil {
		return nil, err
	}

	e, err := (*s.accountRepo).FindBy(ctx, cpf)
	if appErr, ok := err.(*types.Err); ok && appErr.Code == types.EmptyResultErr {
		return nil, types.NewErr(types.AuthenticationErr, "account with the given cpf does not exist", &err)
	}
	if err != nil {
		log.Info().Caller().Err(err).Str("cpf", cpf).Msg("unable to find the account entity via cpf")
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(e.Secret), []byte(secret))
	if err != nil {
		return nil, types.NewErr(types.AuthenticationErr, "the provided secret doesn't match the account's secret", &err)
	}
	return dto.NewAccountView(e), nil
}

// NewAccount returns a value responsible for managing entity.Account actions and integrity
func NewAccount(txr repository.Transactioner, accountRepo repository.Account) Account {
	return &account{
		accountRepo: &accountRepo,
		txr:         &txr,
		accountVali: &validation.Account{
			AccountRepo: &accountRepo,
		},
	}
}
