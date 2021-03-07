package service_test

import (
	"context"
	"testing"

	"github.com/rafael-sousa/stn-accounts/pkg/model/dto"
	"github.com/rafael-sousa/stn-accounts/pkg/model/entity"
	"github.com/rafael-sousa/stn-accounts/pkg/model/types"
	"github.com/rafael-sousa/stn-accounts/pkg/repository"
	"github.com/rafael-sousa/stn-accounts/pkg/service"
)

func TestAccountServiceCreate(t *testing.T) {
	tt := []struct {
		name      string
		d         *dto.AccountCreation
		repo      func(*dto.AccountCreation) repository.Account
		assertErr func(*testing.T, error)
	}{
		{
			name: "create account successfully",
			repo: func(d *dto.AccountCreation) repository.Account {
				return &accountRepoMock{
					findBy: func(ctx context.Context, cpf string) (*entity.Account, error) {
						assertEq(t, "cpf", d.CPF, cpf)
						return nil, nil
					},
					create: func(ctx context.Context, e *entity.Account) (*entity.Account, error) {
						assertEq(t, "name", d.Name, e.Name)
						assertEq(t, "balance", d.Balance, e.Balance.Float64())
						assertEq(t, "cpf", d.CPF, e.CPF)
						e.ID = 1
						return e, nil
					},
				}
			},
			d: newAccountCreation("John", "00000000000", "pw", 100),
		},
		{
			name: "create account with validation err",
			repo: func(d *dto.AccountCreation) repository.Account {
				return &accountRepoMock{}
			},
			d: &dto.AccountCreation{},
			assertErr: func(t *testing.T, err error) {
				assertCustomErr(t, types.ValidationErr, err, "field 'name' is required")

			},
		},
		{
			name: "create account with repository error",
			repo: func(d *dto.AccountCreation) repository.Account {
				return &accountRepoMock{
					findBy: func(ctx context.Context, cpf string) (*entity.Account, error) {
						return nil, nil
					},
					create: func(ctx context.Context, e *entity.Account) (*entity.Account, error) {
						return nil, types.NewErr(types.InternalErr, "internal error", nil)
					},
				}
			},
			d: newAccountCreation("Doe", "98765432100", "l", 0),
			assertErr: func(t *testing.T, err error) {
				assertCustomErr(t, types.InternalErr, err, "internal error")
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s := service.NewAccount(&txr, tc.repo(tc.d))
			acc, err := s.Create(context.Background(), tc.d)
			if err == nil && tc.assertErr == nil {
				assertNotDefault(t, "id", acc.ID)
				assertEq(t, "name", tc.d.Name, acc.Name)
				assertEq(t, "balance", tc.d.Balance, acc.Balance)
				assertEq(t, "cpf", tc.d.CPF, acc.CPF)
				assertNotDefault(t, "created_at", acc.CreatedAt)
				return
			}
			tc.assertErr(t, err)
		})
	}
}

func TestAccountServiceFetch(t *testing.T) {
	tt := []struct {
		name         string
		expectedSize int
		repo         func() repository.Account
		assertErr    func(*testing.T, error)
	}{
		{
			name:         "fetch account with no result",
			expectedSize: 0,
			repo: func() repository.Account {
				return &accountRepoMock{
					fetch: func(ctx context.Context) ([]*entity.Account, error) {
						return []*entity.Account{}, nil
					},
				}
			},
		},
		{
			name:         "fetch account with results",
			expectedSize: 3,
			repo: func() repository.Account {
				return &accountRepoMock{
					fetch: func(ctx context.Context) ([]*entity.Account, error) {
						return []*entity.Account{
							newAccount(1, "Jose", "00123456789", "PW001", 100),
							newAccount(2, "Maria", "98765432100", "PW002", 200),
							newAccount(3, "Silva", "98745632100", "PW003", 300),
						}, nil
					},
				}
			},
		},
		{
			name: "fetch account repository error",
			repo: func() repository.Account {
				return &accountRepoMock{
					fetch: func(ctx context.Context) ([]*entity.Account, error) {
						return nil, types.NewErr(types.InternalErr, "internal error", nil)
					},
				}
			},
			assertErr: func(t *testing.T, err error) {
				assertCustomErr(t, types.InternalErr, err, "internal error")
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s := service.NewAccount(&txr, tc.repo())
			accs, err := s.Fetch(context.Background())
			if err == nil && tc.assertErr == nil {
				assertEq(t, "accs size", len(accs), tc.expectedSize)
				for _, acc := range accs {
					assertNotDefault(t, "name", acc.Name)
					assertNotDefault(t, "cpf", acc.CPF)
					assertNotDefault(t, "balance", acc.Balance)
					assertNotDefault(t, "created_at", acc.CreatedAt)
				}
			} else {
				tc.assertErr(t, err)
			}
		})
	}
}

func TestAccountServiceGetBalance(t *testing.T) {
	tt := []struct {
		name      string
		expected  float64
		repo      func(int64, float64) repository.Account
		assertErr func(*testing.T, error)
		id        int64
	}{
		{
			name: "get account balance successfully",
			repo: func(id int64, balance float64) repository.Account {
				return &accountRepoMock{
					getBalance: func(ctx context.Context, currentID int64) (types.Currency, error) {
						assertEq(t, "id", id, currentID)
						return types.NewCurrency(balance), nil
					},
				}
			},
			expected: 500,
			id:       1,
		},
		{
			name: "get account balance with repository error",
			repo: func(id int64, balance float64) repository.Account {
				return &accountRepoMock{
					getBalance: func(ctx context.Context, currentID int64) (types.Currency, error) {
						return types.NewCurrency(0), types.NewErr(types.EmptyResultErr, "no sql rows", nil)
					},
				}
			},
			expected: 0,
			id:       2,
			assertErr: func(t *testing.T, err error) {
				assertCustomErr(t, types.EmptyResultErr, err, "no sql rows")
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s := service.NewAccount(&txr, tc.repo(tc.id, tc.expected))
			balance, err := s.GetBalance(context.Background(), tc.id)
			assertEq(t, "balance", tc.expected, balance)
			if err != nil {
				tc.assertErr(t, err)
			}
		})
	}
}

func TestAccountServiceLogin(t *testing.T) {
	tt := []struct {
		name      string
		expected  *entity.Account
		repo      func(*entity.Account) repository.Account
		assertErr func(*testing.T, error)
		cpf       string
		secret    string
	}{
		{
			name:     "login with cpf and secret successfully",
			expected: newAccount(1, "Sousa", "11111111111", "$2a$10$c3GzxvPAAMS9pDqB9XIYi.kT/PN7CxfRev.BsRLvAJqVcZnFiW05i", 0),
			secret:   "...",
			cpf:      "11111111111",
			repo: func(exp *entity.Account) repository.Account {
				return &accountRepoMock{
					findBy: func(ctx context.Context, cpf string) (*entity.Account, error) {
						assertEq(t, "cpf", exp.CPF, cpf)
						return exp, nil
					},
				}
			},
		},
		{
			name:     "login with cpf and wrong secret",
			expected: newAccount(2, "Alice", "22222222222", "123", 10),
			secret:   "123",
			cpf:      "22222222222",
			repo: func(exp *entity.Account) repository.Account {
				return &accountRepoMock{
					findBy: func(ctx context.Context, cpf string) (*entity.Account, error) {
						return exp, nil
					},
				}
			},
			assertErr: func(t *testing.T, err error) {
				assertCustomErr(t, types.AuthenticationErr, err, "the provided secret doesn't match the account's secret")
			},
		},
		{
			name:   "login with repository error",
			secret: "...",
			cpf:    "33333333333",
			repo: func(exp *entity.Account) repository.Account {
				return &accountRepoMock{
					findBy: func(ctx context.Context, cpf string) (*entity.Account, error) {
						return exp, types.NewErr(types.InternalErr, "internal error", nil)
					},
				}
			},
			assertErr: func(t *testing.T, err error) {
				assertCustomErr(t, types.InternalErr, err, "internal error")
			},
		},
		{
			name:   "login with nonexistent account",
			secret: "...",
			cpf:    "44444444444",
			repo: func(exp *entity.Account) repository.Account {
				return &accountRepoMock{
					findBy: func(ctx context.Context, cpf string) (*entity.Account, error) {
						return exp, types.NewErr(types.EmptyResultErr, "no result", nil)
					},
				}
			},
			assertErr: func(t *testing.T, err error) {
				assertCustomErr(t, types.AuthenticationErr, err, "account with the given cpf does not exist")
			},
		},
		{
			name:   "login without cpf",
			secret: "...",
			repo: func(exp *entity.Account) repository.Account {
				return &accountRepoMock{}
			},
			assertErr: func(t *testing.T, err error) {
				assertCustomErr(t, types.ValidationErr, err, "field 'cpf' is required")
			},
		},
		{
			name: "login without secret",
			cpf:  "55555555555",
			repo: func(exp *entity.Account) repository.Account {
				return &accountRepoMock{}
			},
			assertErr: func(t *testing.T, err error) {
				assertCustomErr(t, types.ValidationErr, err, "field 'secret' is required")
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s := service.NewAccount(&txr, tc.repo(tc.expected))
			view, err := s.Login(context.Background(), tc.cpf, tc.secret)
			if err != nil {
				tc.assertErr(t, err)
			} else {
				assertEq(t, "id", tc.expected.ID, view.ID)
				assertEq(t, "name", tc.expected.Name, view.Name)
				assertEq(t, "cpf", tc.expected.CPF, view.CPF)
				assertEq(t, "balance", tc.expected.Balance.Float64(), view.Balance)
				assertNotDefault(t, "created_at", view.CreatedAt)
			}
		})
	}
}
