package validation_test

import (
	"context"
	"testing"

	"github.com/rafael-sousa/stn-accounts/pkg/model/dto"
	"github.com/rafael-sousa/stn-accounts/pkg/model/entity"
	"github.com/rafael-sousa/stn-accounts/pkg/model/types"
	"github.com/rafael-sousa/stn-accounts/pkg/repository"
	"github.com/rafael-sousa/stn-accounts/pkg/service/validation"
)

func TestAccountCreation(t *testing.T) {
	tt := []struct {
		name            string
		repo            func() repository.Account
		accountCreation *dto.AccountCreation
		assertErr       func(*testing.T, error)
	}{
		{
			name: "validate account creation successfully",
			repo: func() repository.Account {
				return &accountRepoMock{
					findBy: func(c context.Context, cpf string) (*entity.Account, error) {
						return nil, nil
					},
				}
			},
			assertErr: func(t *testing.T, err error) {
				if err != nil {
					t.Errorf("expected no error but got %v", err)
				}
			},
			accountCreation: newAccountCreation("John", "12345678900", "pw", 999),
		},
		{
			name: "validate account creation with empty name",
			repo: func() repository.Account {
				return &accountRepoMock{}
			},
			assertErr: func(t *testing.T, err error) {
				assertCustomErr(t, types.ValidationErr, err, "field 'name' is required")
			},
			accountCreation: newAccountCreation("", "12345678900", "pw", 999),
		},
		{
			name: "validate account creation with name having trailing white space",
			repo: func() repository.Account {
				return &accountRepoMock{}
			},
			assertErr: func(t *testing.T, err error) {
				assertCustomErr(t, types.ValidationErr, err, "field 'name' can't have trailing whitespace")
			},
			accountCreation: newAccountCreation(" Dan", "12345678900", "pw", 999),
		},
		{
			name: "validate account creation with cpf length less than 11 chars",
			repo: func() repository.Account {
				return &accountRepoMock{}
			},
			assertErr: func(t *testing.T, err error) {
				assertCustomErr(t, types.ValidationErr, err, "field 'cpf' must have at least 11 characters")
			},
			accountCreation: newAccountCreation("Jack", "2345678900", "pw", 999),
		},
		{
			name: "validate account creation with cpf length more than 11 chars",
			repo: func() repository.Account {
				return &accountRepoMock{}
			},
			assertErr: func(t *testing.T, err error) {
				assertCustomErr(t, types.ValidationErr, err, "field 'cpf' must have at most 11 characters")
			},
			accountCreation: newAccountCreation("Bia", "123456789000", "pw", 999),
		},
		{
			name: "validate account creation with empty cpf",
			repo: func() repository.Account {
				return &accountRepoMock{}
			},
			assertErr: func(t *testing.T, err error) {
				assertCustomErr(t, types.ValidationErr, err, "field 'cpf' is required")
			},
			accountCreation: newAccountCreation("Teo", "", "pw", 999),
		},
		{
			name: "validate account creation with invalid cpf",
			repo: func() repository.Account {
				return &accountRepoMock{}
			},
			assertErr: func(t *testing.T, err error) {
				assertCustomErr(t, types.ValidationErr, err, "field 'cpf' has an invalid format")
			},
			accountCreation: newAccountCreation("Elen", "0000000000#", "pw", 999),
		},
		{
			name: "validate account creation with empty secret",
			repo: func() repository.Account {
				return &accountRepoMock{}
			},
			assertErr: func(t *testing.T, err error) {
				assertCustomErr(t, types.ValidationErr, err, "field 'secret' is required")
			},
			accountCreation: newAccountCreation("Olly", "00000000000", "", 999),
		},
		{
			name: "validate account creation with negative balance",
			repo: func() repository.Account {
				return &accountRepoMock{}
			},
			assertErr: func(t *testing.T, err error) {
				assertCustomErr(t, types.ValidationErr, err, "field 'balance' must be greater than or equal to 0")
			},
			accountCreation: newAccountCreation("Paul", "11111111111", "pw", -0.01),
		},
		{
			name: "validate account creation with existing cpf",
			repo: func() repository.Account {
				return &accountRepoMock{
					findBy: func(c context.Context, s string) (*entity.Account, error) {
						return &entity.Account{}, nil
					},
				}
			},
			assertErr: func(t *testing.T, err error) {
				assertCustomErr(t, types.ConflictErr, err, "field 'cpf' with value '22222222222' is already in use")
			},
			accountCreation: newAccountCreation("Carol", "22222222222", "pw", 1),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo := tc.repo()
			v := validation.Account{
				AccountRepository: &repo,
			}
			err := v.Creation(context.Background(), tc.accountCreation)
			tc.assertErr(t, err)
		})
	}
}

func TestAccountLogin(t *testing.T) {
	tt := []struct {
		name      string
		cpf       string
		secret    string
		assertErr func(*testing.T, error)
	}{
		{
			name:   "validate login creation successfully",
			cpf:    "00000000000",
			secret: "pw",
			assertErr: func(t *testing.T, err error) {
				if err != nil {
					t.Errorf("expected no error but got %v", err)
				}
			},
		},
		{
			name:   "validate login creation with no cpf",
			cpf:    "",
			secret: "pw",
			assertErr: func(t *testing.T, err error) {
				assertCustomErr(t, types.ValidationErr, err, "field 'cpf' is required")
			},
		},
		{
			name:   "validate login creation with no secret",
			cpf:    "00000000000",
			secret: "",
			assertErr: func(t *testing.T, err error) {
				assertCustomErr(t, types.ValidationErr, err, "field 'secret' is required")
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			v := validation.Account{}
			err := v.Login(tc.cpf, tc.secret)
			tc.assertErr(t, err)
		})
	}
}
