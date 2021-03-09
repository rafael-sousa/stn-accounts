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
	accountRepository *repository.Account
	accountValidator  *validation.Account
	txr               *repository.Transactioner
}

// Fetch returns a list of dto.AccountView
func (srv *account) Fetch(ctx context.Context) ([]*dto.AccountView, error) {
	accounts, err := (*srv.accountRepository).Fetch(ctx)
	if err != nil {
		log.Error().Caller().Err(err).Msg("unable to fetch accounts")
		return nil, err
	}
	views := make([]*dto.AccountView, 0, len(accounts))
	for _, account := range accounts {
		views = append(views, dto.NewAccountView(account))
	}
	return views, nil
}

// GetBalance returns the given account balance
func (srv *account) GetBalance(ctx context.Context, id int64) (float64, error) {
	balance, err := (*srv.accountRepository).GetBalance(ctx, id)
	if err != nil {
		log.Info().Caller().Err(err).Int64("id", id).Msg("unable to get account balance")
		return 0, err
	}
	return balance.Float64(), nil
}

// Create validates and persists the given e entity.Account
func (srv *account) Create(ctx context.Context, accountCreation *dto.AccountCreation) (*dto.AccountView, error) {
	var account *entity.Account
	err := (*srv.txr).WithTx(ctx, func(txCtx context.Context) error {
		err := srv.accountValidator.Creation(txCtx, accountCreation)
		if err != nil {
			return err
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(accountCreation.Secret), bcrypt.DefaultCost)
		if err != nil {
			log.Info().Caller().Err(err).Msg("unable to create the account secret hash")
			return err
		}
		account, err = (*srv.accountRepository).Create(txCtx, &entity.Account{
			Name:      accountCreation.Name,
			CPF:       accountCreation.CPF,
			Balance:   types.NewCurrency(accountCreation.Balance),
			CreatedAt: time.Now(),
			Secret:    string(hash),
		})
		return err
	})

	if err != nil {
		log.Info().Caller().Err(err).Str("cpf", accountCreation.CPF).Str("name", accountCreation.Name).Msg("unable to create account")
		return nil, err
	}

	return dto.NewAccountView(account), nil
}

// Login is responsible for fetching an account, comparing the secret, and returning the respective account view.
// It returns nil and an err if the account is not found or the provided secret doesn't match the account
func (srv *account) Login(ctx context.Context, cpf string, secret string) (*dto.AccountView, error) {
	err := srv.accountValidator.Login(cpf, secret)
	if err != nil {
		return nil, err
	}

	account, err := (*srv.accountRepository).FindBy(ctx, cpf)
	if customErr, ok := err.(*types.Err); ok && customErr.Code == types.EmptyResultErr {
		return nil, types.NewErr(types.AuthenticationErr, "account with the given cpf does not exist", &err)
	}
	if err != nil {
		log.Info().Caller().Err(err).Str("cpf", cpf).Msg("unable to find the account entity via cpf")
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Secret), []byte(secret))
	if err != nil {
		return nil, types.NewErr(types.AuthenticationErr, "the provided secret doesn't match the account's secret", &err)
	}
	return dto.NewAccountView(account), nil
}

// NewAccount returns a value responsible for managing entity.Account actions and integrity
func NewAccount(txr *repository.Transactioner, accountRepository *repository.Account) Account {
	return &account{
		accountRepository: accountRepository,
		txr:               txr,
		accountValidator: &validation.Account{
			AccountRepository: accountRepository,
		},
	}
}
