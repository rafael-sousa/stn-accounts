package validation_test

import (
	"context"
	"testing"

	"github.com/rafael-sousa/stn-accounts/pkg/model/dto"
	"github.com/rafael-sousa/stn-accounts/pkg/model/types"
	"github.com/rafael-sousa/stn-accounts/pkg/repository"
	"github.com/rafael-sousa/stn-accounts/pkg/service/validation"
	"github.com/rafael-sousa/stn-accounts/pkg/testutil"
)

func TestTransferCreation(t *testing.T) {
	tt := []struct {
		name             string
		repo             func() repository.Account
		transferCreation *dto.TransferCreation
		origin           int64
		assertErr        func(*testing.T, error)
	}{
		{
			name: "validate transfer creation successfully",
			repo: func() repository.Account {
				return &testutil.AccountRepoMock{
					ExpectGetBalance: func(c context.Context, i int64) (types.Currency, error) {
						testutil.AssertEq(t, "origin id", int64(1), i)
						return types.NewCurrency(1000), nil
					},
					ExpectExists: func(c context.Context, i int64) (bool, error) {
						testutil.AssertEq(t, "destination id", int64(2), i)
						return true, nil
					},
				}
			},
			assertErr: testutil.AssertNoErr,
			origin:    1,
			transferCreation: &dto.TransferCreation{
				Destination: 2,
				Amount:      500,
			},
		},
		{
			name: "validate transfer creation with destination equal origin",
			repo: func() repository.Account {
				return &testutil.AccountRepoMock{}
			},
			assertErr: func(t *testing.T, err error) {
				testutil.AssertCustomErr(t, types.ConflictErr, err, "fields 'origin id' and 'destination id' can't be the same")
			},
			origin: 1,
			transferCreation: &dto.TransferCreation{
				Destination: 1,
				Amount:      500,
			},
		},
		{
			name: "validate transfer creation with no amount",
			repo: func() repository.Account {
				return &testutil.AccountRepoMock{}
			},
			assertErr: func(t *testing.T, err error) {
				testutil.AssertCustomErr(t, types.ValidationErr, err, "field 'amount' must be greater than 0")
			},
			origin: 1,
			transferCreation: &dto.TransferCreation{
				Destination: 2,
				Amount:      0,
			},
		},
		{
			name: "validate transfer creation with no destination",
			repo: func() repository.Account {
				return &testutil.AccountRepoMock{}
			},
			assertErr: func(t *testing.T, err error) {
				testutil.AssertCustomErr(t, types.ValidationErr, err, "field 'destination_id' is required")
			},
			origin: 1,
			transferCreation: &dto.TransferCreation{
				Destination: 0,
				Amount:      0.01,
			},
		},
		{
			name: "validate transfer creation with no funds",
			repo: func() repository.Account {
				return &testutil.AccountRepoMock{
					ExpectGetBalance: func(c context.Context, i int64) (types.Currency, error) {
						return 0, nil
					},
				}
			},
			assertErr: func(t *testing.T, err error) {
				testutil.AssertCustomErr(t, types.ValidationErr, err, "the origin must have a balance greater than or equal to 50.00")
			},
			origin: 1,
			transferCreation: &dto.TransferCreation{
				Destination: 3,
				Amount:      50,
			},
		},
		{
			name: "validate transfer creation from non existent origin",
			repo: func() repository.Account {
				return &testutil.AccountRepoMock{
					ExpectGetBalance: func(c context.Context, i int64) (types.Currency, error) {
						return 0, types.NewErr(types.EmptyResultErr, "no row return", nil)
					},
				}
			},
			assertErr: func(t *testing.T, err error) {
				testutil.AssertCustomErr(t, types.EmptyResultErr, err, "record with 'origin' equals '1' was not found")
			},
			origin: 1,
			transferCreation: &dto.TransferCreation{
				Destination: 2,
				Amount:      10,
			},
		},
		{
			name: "validate transfer creation to non existent destination",
			repo: func() repository.Account {
				return &testutil.AccountRepoMock{
					ExpectGetBalance: func(c context.Context, i int64) (types.Currency, error) {
						return types.NewCurrency(500), nil
					},
					ExpectExists: func(c context.Context, i int64) (bool, error) {
						testutil.AssertEq(t, "destination id", int64(3), i)
						return false, nil
					},
				}
			},
			assertErr: func(t *testing.T, err error) {
				testutil.AssertCustomErr(t, types.EmptyResultErr, err, "record with 'destination' equals '3' was not found")
			},
			origin: 1,
			transferCreation: &dto.TransferCreation{
				Destination: 3,
				Amount:      10,
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo := tc.repo()
			v := validation.Transfer{
				AccountRepository: &repo,
			}
			err := v.Creation(context.Background(), tc.origin, *tc.transferCreation)
			tc.assertErr(t, err)
		})
	}
}
