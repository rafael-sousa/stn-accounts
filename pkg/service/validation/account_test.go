package validation_test

import (
	"context"
	"testing"

	"github.com/rafael-sousa/stn-accounts/pkg/model/dto"
	"github.com/rafael-sousa/stn-accounts/pkg/model/entity"
	"github.com/rafael-sousa/stn-accounts/pkg/model/types"
	"github.com/rafael-sousa/stn-accounts/pkg/repository"
	"github.com/rafael-sousa/stn-accounts/pkg/service/validation"
	"github.com/rafael-sousa/stn-accounts/pkg/testutil"
)

func TestAccountCreation(t *testing.T) {
	tt := []struct {
		name            string
		repo            func() repository.Account
		accountCreation dto.AccountCreation
		assertErr       func(*testing.T, error)
	}{
		{
			name: "validate account creation successfully",
			repo: func() repository.Account {
				return &testutil.AccountRepoMock{
					ExpectFindBy: func(c context.Context, cpf string) (entity.Account, error) {
						return entity.Account{}, types.NewErr(types.EmptyResultErr, "EmptyResultErr", nil)
					},
				}
			},
			assertErr:       testutil.AssertNoErr,
			accountCreation: testutil.NewAccountCreation("John", "71453945024", "pw", 999),
		},
		{
			name: "validate account creation with empty name",
			repo: func() repository.Account {
				return &testutil.AccountRepoMock{}
			},
			assertErr: func(t *testing.T, err error) {
				testutil.AssertCustomErr(t, types.ValidationErr, err, "field 'name' is required")
			},
			accountCreation: testutil.NewAccountCreation("", "12345678900", "pw", 999),
		},
		{
			name: "validate account creation with name having trailing white space",
			repo: func() repository.Account {
				return &testutil.AccountRepoMock{}
			},
			assertErr: func(t *testing.T, err error) {
				testutil.AssertCustomErr(t, types.ValidationErr, err, "field 'name' can't have trailing whitespace")
			},
			accountCreation: testutil.NewAccountCreation(" Dan", "71453945024", "pw", 999),
		},
		{
			name: "validate account creation with empty cpf",
			repo: func() repository.Account {
				return &testutil.AccountRepoMock{}
			},
			assertErr: func(t *testing.T, err error) {
				testutil.AssertCustomErr(t, types.ValidationErr, err, "field 'cpf' is required")
			},
			accountCreation: testutil.NewAccountCreation("Teo", "", "pw", 999),
		},
		{
			name: "validate account creation with invalid cpf",
			repo: func() repository.Account {
				return &testutil.AccountRepoMock{}
			},
			assertErr: func(t *testing.T, err error) {
				testutil.AssertCustomErr(t, types.ValidationErr, err, "field 'cpf' has an invalid format")
			},
			accountCreation: testutil.NewAccountCreation("Elen", "0000000000#", "pw", 999),
		},
		{
			name: "validate account creation with empty secret",
			repo: func() repository.Account {
				return &testutil.AccountRepoMock{}
			},
			assertErr: func(t *testing.T, err error) {
				testutil.AssertCustomErr(t, types.ValidationErr, err, "field 'secret' is required")
			},
			accountCreation: testutil.NewAccountCreation("Olly", "18474605008", "", 999),
		},
		{
			name: "validate account creation with negative balance",
			repo: func() repository.Account {
				return &testutil.AccountRepoMock{}
			},
			assertErr: func(t *testing.T, err error) {
				testutil.AssertCustomErr(t, types.ValidationErr, err, "field 'balance' must be greater than or equal to 0")
			},
			accountCreation: testutil.NewAccountCreation("Paul", "44206294011", "pw", -0.01),
		},
		{
			name: "validate account creation with existing cpf",
			repo: func() repository.Account {
				return &testutil.AccountRepoMock{
					ExpectFindBy: func(c context.Context, s string) (entity.Account, error) {
						return entity.Account{}, nil
					},
				}
			},
			assertErr: func(t *testing.T, err error) {
				testutil.AssertCustomErr(t, types.ConflictErr, err, "field 'cpf' with value '28494002031' is already in use")
			},
			accountCreation: testutil.NewAccountCreation("Carol", "28494002031", "pw", 1),
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
			name:      "validate login creation successfully",
			cpf:       "06928541008",
			secret:    "pw",
			assertErr: testutil.AssertNoErr,
		},
		{
			name:   "validate login creation with no cpf",
			cpf:    "",
			secret: "pw",
			assertErr: func(t *testing.T, err error) {
				testutil.AssertCustomErr(t, types.ValidationErr, err, "field 'cpf' is required")
			},
		},
		{
			name:   "validate login creation with no secret",
			cpf:    "52803492083",
			secret: "",
			assertErr: func(t *testing.T, err error) {
				testutil.AssertCustomErr(t, types.ValidationErr, err, "field 'secret' is required")
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
