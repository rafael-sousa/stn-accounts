package service_test

import (
	"context"
	"testing"

	"github.com/rafael-sousa/stn-accounts/pkg/model/dto"
	"github.com/rafael-sousa/stn-accounts/pkg/model/entity"
	"github.com/rafael-sousa/stn-accounts/pkg/model/types"
	"github.com/rafael-sousa/stn-accounts/pkg/repository"
	"github.com/rafael-sousa/stn-accounts/pkg/service"
	"github.com/rafael-sousa/stn-accounts/pkg/testutil"
)

func TestTransferServiceFetch(t *testing.T) {
	tt := []struct {
		name         string
		expectedSize int
		transferRepo func(int64) repository.Transfer
		accountRepo  func() repository.Account
		assertErr    func(*testing.T, error)
		id           int64
	}{
		{
			name:         "fetch transfers with no result",
			expectedSize: 0,
			transferRepo: func(id int64) repository.Transfer {
				return &testutil.TransferRepoMock{
					ExpectFetch: func(ctx context.Context, currentID int64) ([]*entity.Transfer, error) {
						testutil.AssertEq(t, "id", id, currentID)
						return []*entity.Transfer{}, nil
					},
				}
			},
			accountRepo: func() repository.Account {
				return &testutil.AccountRepoMock{}
			},
			id: 1,
		},
		{
			name:         "fetch transfers with repository error",
			expectedSize: 0,
			transferRepo: func(id int64) repository.Transfer {
				return &testutil.TransferRepoMock{
					ExpectFetch: func(ctx context.Context, currentID int64) ([]*entity.Transfer, error) {
						return nil, types.NewErr(types.InternalErr, "internal error", nil)
					},
				}
			},
			accountRepo: func() repository.Account {
				return &testutil.AccountRepoMock{}
			},
			id: 2,
			assertErr: func(t *testing.T, err error) {
				testutil.AssertCustomErr(t, types.InternalErr, err, "internal error")
			},
		},
		{
			name:         "fetch transfers successfully with results",
			expectedSize: 5,
			transferRepo: func(id int64) repository.Transfer {
				return &testutil.TransferRepoMock{
					ExpectFetch: func(ctx context.Context, currentID int64) ([]*entity.Transfer, error) {
						testutil.AssertEq(t, "id", id, currentID)
						return []*entity.Transfer{
							testutil.NewEntityTransfer(1, 3, 2, 10),
							testutil.NewEntityTransfer(2, 3, 2, 20),
							testutil.NewEntityTransfer(3, 3, 2, 30),
							testutil.NewEntityTransfer(4, 3, 3, 40),
							testutil.NewEntityTransfer(5, 3, 3, 50),
						}, nil
					},
				}
			},
			accountRepo: func() repository.Account {
				return &testutil.AccountRepoMock{}
			},
			id: 3,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			transferRepo := tc.transferRepo(tc.id)
			accRepo := tc.accountRepo()
			s := service.NewTransfer(&txr, &transferRepo, &accRepo)
			transfers, err := s.Fetch(context.Background(), tc.id)
			if err == nil && tc.assertErr == nil {
				testutil.AssertEq(t, "transfers size", len(transfers), tc.expectedSize)
				for _, transfer := range transfers {
					testutil.AssertNotDefault(t, "amount", transfer.Amount)
				}
			} else {
				tc.assertErr(t, err)
			}
		})
	}
}

func TestTransferServiceCreate(t *testing.T) {
	tt := []struct {
		name         string
		expected     *dto.TransferView
		transferRepo func(int64, *dto.TransferCreation) repository.Transfer
		accountRepo  func(int64, *dto.TransferCreation) repository.Account
		assertErr    func(*testing.T, error)
		origin       int64
		d            *dto.TransferCreation
	}{
		{
			name: "create transfer successfully",
			transferRepo: func(origin int64, d *dto.TransferCreation) repository.Transfer {
				return &testutil.TransferRepoMock{
					ExpectCreate: func(ctx context.Context, e *entity.Transfer) (*entity.Transfer, error) {
						testutil.AssertEq(t, "origin", origin, e.Origin)
						testutil.AssertEq(t, "destination", d.Destination, e.Destination)
						testutil.AssertEq(t, "amount", types.NewCurrency(d.Amount), e.Amount)
						testutil.AssertNotDefault(t, "created_at", e.CreatedAt)
						e.ID = 1
						return e, nil
					},
				}
			},
			accountRepo: func(origin int64, d *dto.TransferCreation) repository.Account {
				balanceStack := []float64{500, 500, 0}
				return &testutil.AccountRepoMock{
					ExpectGetBalance: func(c context.Context, i int64) (types.Currency, error) {
						if len(balanceStack) == 0 {
							t.Fatalf("unexpected method call")
						}
						b := types.NewCurrency(balanceStack[0])
						balanceStack = balanceStack[1:]
						return b, nil
					},
					ExpectUpdateBalance: func(ctx context.Context, i int64, b types.Currency) error {
						switch i {
						case origin:
							testutil.AssertEq(t, "origin balance", types.NewCurrency(0), b)
						case d.Destination:
							testutil.AssertEq(t, "destination balance", types.NewCurrency(500), b)
						default:
							t.Fatalf("unexpected method call")
						}
						return nil
					},
					ExpectExists: func(c context.Context, i int64) (bool, error) {
						testutil.AssertEq(t, "destination id", d.Destination, i)
						return true, nil
					},
				}
			},
			origin: 1,
			d: &dto.TransferCreation{
				Destination: 2,
				Amount:      500,
			},
		},
		{
			name: "create transfer with insufficient funds",
			transferRepo: func(origin int64, d *dto.TransferCreation) repository.Transfer {
				return &testutil.TransferRepoMock{}
			},
			accountRepo: func(origin int64, d *dto.TransferCreation) repository.Account {
				return &testutil.AccountRepoMock{
					ExpectGetBalance: func(c context.Context, i int64) (types.Currency, error) {
						return types.NewCurrency(0), nil
					},
				}
			},
			origin: 1,
			d: &dto.TransferCreation{
				Destination: 2,
				Amount:      500,
			},
			assertErr: func(t *testing.T, err error) {
				testutil.AssertCustomErr(t, types.ValidationErr, err, "the origin must have a balance greater than or equal to 500.00")
			},
		},
		{
			name: "create transfer with destination equal to origin",
			transferRepo: func(origin int64, d *dto.TransferCreation) repository.Transfer {
				return &testutil.TransferRepoMock{}
			},
			accountRepo: func(origin int64, d *dto.TransferCreation) repository.Account {
				return &testutil.AccountRepoMock{}
			},
			origin: 1,
			d: &dto.TransferCreation{
				Destination: 1,
				Amount:      100,
			},
			assertErr: func(t *testing.T, err error) {
				testutil.AssertCustomErr(t, types.ConflictErr, err, "fields 'origin id' and 'destination id' can't be the same")
			},
		},
		{
			name: "create transfer with repository error getting origin balance",
			transferRepo: func(origin int64, d *dto.TransferCreation) repository.Transfer {
				return &testutil.TransferRepoMock{}
			},
			accountRepo: func(origin int64, d *dto.TransferCreation) repository.Account {
				balanceStack := []float64{500}
				return &testutil.AccountRepoMock{
					ExpectGetBalance: func(c context.Context, i int64) (types.Currency, error) {
						if len(balanceStack) == 0 {
							return 0, types.NewErr(types.InternalErr, "internal error", nil)
						}
						b := types.NewCurrency(balanceStack[0])
						balanceStack = balanceStack[1:]
						return b, nil
					},
					ExpectExists: func(c context.Context, i int64) (bool, error) {
						testutil.AssertEq(t, "destination id", d.Destination, i)
						return true, nil
					},
				}
			},
			origin: 1,
			d: &dto.TransferCreation{
				Destination: 2,
				Amount:      100,
			},
			assertErr: func(t *testing.T, err error) {
				testutil.AssertCustomErr(t, types.InternalErr, err, "internal error")
			},
		},
		{
			name: "create transfer with repository error getting destination balance",
			transferRepo: func(origin int64, d *dto.TransferCreation) repository.Transfer {
				return &testutil.TransferRepoMock{}
			},
			accountRepo: func(origin int64, d *dto.TransferCreation) repository.Account {
				balanceStack := []float64{500, 500}
				return &testutil.AccountRepoMock{
					ExpectGetBalance: func(c context.Context, i int64) (types.Currency, error) {
						if len(balanceStack) == 0 {
							return 0, types.NewErr(types.InternalErr, "internal error", nil)
						}
						b := types.NewCurrency(balanceStack[0])
						balanceStack = balanceStack[1:]
						return b, nil
					},
					ExpectExists: func(c context.Context, i int64) (bool, error) {
						testutil.AssertEq(t, "destination id", d.Destination, i)
						return true, nil
					},
				}
			},
			origin: 1,
			d: &dto.TransferCreation{
				Destination: 2,
				Amount:      100,
			},
			assertErr: func(t *testing.T, err error) {
				testutil.AssertCustomErr(t, types.InternalErr, err, "internal error")
			},
		},
		{
			name: "create transfer with repository error updating origin balance",
			transferRepo: func(origin int64, d *dto.TransferCreation) repository.Transfer {
				return &testutil.TransferRepoMock{}
			},
			accountRepo: func(origin int64, d *dto.TransferCreation) repository.Account {
				balanceStack := []float64{500, 500, 0}
				return &testutil.AccountRepoMock{
					ExpectGetBalance: func(c context.Context, i int64) (types.Currency, error) {
						if len(balanceStack) == 0 {
							return 0, types.NewErr(types.InternalErr, "internal error", nil)
						}
						b := types.NewCurrency(balanceStack[0])
						balanceStack = balanceStack[1:]
						return b, nil
					},
					ExpectUpdateBalance: func(ctx context.Context, i int64, b types.Currency) error {
						if i == origin {
							return types.NewErr(types.InternalErr, "internal error", nil)
						}
						t.Fatalf("unexpected method call")
						return nil
					},
					ExpectExists: func(c context.Context, i int64) (bool, error) {
						testutil.AssertEq(t, "destination id", d.Destination, i)
						return true, nil
					},
				}
			},
			origin: 1,
			d: &dto.TransferCreation{
				Destination: 2,
				Amount:      100,
			},
			assertErr: func(t *testing.T, err error) {
				testutil.AssertCustomErr(t, types.InternalErr, err, "internal error")
			},
		},
		{
			name: "create transfer with repository error updating destination balance",
			transferRepo: func(origin int64, d *dto.TransferCreation) repository.Transfer {
				return &testutil.TransferRepoMock{}
			},
			accountRepo: func(origin int64, d *dto.TransferCreation) repository.Account {
				balanceStack := []float64{500, 500, 0}
				return &testutil.AccountRepoMock{
					ExpectGetBalance: func(c context.Context, i int64) (types.Currency, error) {
						if len(balanceStack) == 0 {
							return 0, types.NewErr(types.InternalErr, "internal error", nil)
						}
						b := types.NewCurrency(balanceStack[0])
						balanceStack = balanceStack[1:]
						return b, nil
					},
					ExpectUpdateBalance: func(ctx context.Context, i int64, b types.Currency) error {
						switch i {
						case origin:
							return nil
						case d.Destination:
							return types.NewErr(types.InternalErr, "internal error", nil)
						default:
							t.Fatalf("unexpected method call")
						}
						return nil
					},
					ExpectExists: func(c context.Context, i int64) (bool, error) {
						testutil.AssertEq(t, "destination id", d.Destination, i)
						return true, nil
					},
				}
			},
			origin: 1,
			d: &dto.TransferCreation{
				Destination: 2,
				Amount:      100,
			},
			assertErr: func(t *testing.T, err error) {
				testutil.AssertCustomErr(t, types.InternalErr, err, "internal error")
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			transferRepo := tc.transferRepo(tc.origin, tc.d)
			accRepo := tc.accountRepo(tc.origin, tc.d)
			s := service.NewTransfer(&txr, &transferRepo, &accRepo)
			view, err := s.Create(context.Background(), tc.origin, tc.d)
			if err == nil && tc.assertErr == nil {
				testutil.AssertNotDefault(t, "id", view.ID)
				testutil.AssertEq(t, "amount", tc.d.Amount, view.Amount)
				testutil.AssertEq(t, "destination", tc.d.Destination, view.Destination)
				testutil.AssertNotDefault(t, "created_at", view.CreatedAt)
			} else {
				tc.assertErr(t, err)
			}
		})
	}
}
