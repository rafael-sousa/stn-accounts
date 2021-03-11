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
				return &testutil.AccountRepoMock{
					ExpectFindBy: func(ctx context.Context, cpf string) (*entity.Account, error) {
						testutil.AssertEq(t, "cpf", d.CPF, cpf)
						return nil, nil
					},
					ExpectCreate: func(ctx context.Context, e *entity.Account) (*entity.Account, error) {
						testutil.AssertEq(t, "name", d.Name, e.Name)
						testutil.AssertEq(t, "balance", d.Balance, e.Balance.Float64())
						testutil.AssertEq(t, "cpf", d.CPF, e.CPF)
						e.ID = 1
						return e, nil
					},
				}
			},
			d:         testutil.NewAccountCreation("John", "62202136029", "pw", 100),
			assertErr: testutil.AssertNoErr,
		},
		{
			name: "create account with validation err",
			repo: func(d *dto.AccountCreation) repository.Account {
				return &testutil.AccountRepoMock{}
			},
			d: &dto.AccountCreation{},
			assertErr: func(t *testing.T, err error) {
				testutil.AssertCustomErr(t, types.ValidationErr, err, "field 'name' is required")
			},
		},
		{
			name: "create account with repository error",
			repo: func(d *dto.AccountCreation) repository.Account {
				return &testutil.AccountRepoMock{
					ExpectFindBy: func(ctx context.Context, cpf string) (*entity.Account, error) {
						return nil, nil
					},
					ExpectCreate: func(ctx context.Context, e *entity.Account) (*entity.Account, error) {
						return nil, types.NewErr(types.InternalErr, "internal error", nil)
					},
				}
			},
			d: testutil.NewAccountCreation("Doe", "98765432100", "l", 0),
			assertErr: func(t *testing.T, err error) {
				testutil.AssertCustomErr(t, types.InternalErr, err, "internal error")
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo := tc.repo(tc.d)
			s := service.NewAccount(&txr, &repo)
			acc, err := s.Create(context.Background(), tc.d)
			if err == nil && tc.assertErr == nil {
				testutil.AssertNotDefault(t, "id", acc.ID)
				testutil.AssertEq(t, "name", tc.d.Name, acc.Name)
				testutil.AssertEq(t, "balance", tc.d.Balance, acc.Balance)
				testutil.AssertEq(t, "cpf", tc.d.CPF, acc.CPF)
				testutil.AssertNotDefault(t, "created_at", acc.CreatedAt)
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
				return &testutil.AccountRepoMock{
					ExpectFetch: func(ctx context.Context) ([]*entity.Account, error) {
						return []*entity.Account{}, nil
					},
				}
			},
		},
		{
			name:         "fetch account with results",
			expectedSize: 3,
			repo: func() repository.Account {
				return &testutil.AccountRepoMock{
					ExpectFetch: func(ctx context.Context) ([]*entity.Account, error) {
						return []*entity.Account{
							testutil.NewEntityAccount(1, "Jose", "00123456789", "PW001", 100),
							testutil.NewEntityAccount(2, "Maria", "98765432100", "PW002", 200),
							testutil.NewEntityAccount(3, "Silva", "98745632100", "PW003", 300),
						}, nil
					},
				}
			},
		},
		{
			name: "fetch account repository error",
			repo: func() repository.Account {
				return &testutil.AccountRepoMock{
					ExpectFetch: func(ctx context.Context) ([]*entity.Account, error) {
						return nil, types.NewErr(types.InternalErr, "internal error", nil)
					},
				}
			},
			assertErr: func(t *testing.T, err error) {
				testutil.AssertCustomErr(t, types.InternalErr, err, "internal error")
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo := tc.repo()
			s := service.NewAccount(&txr, &repo)
			accs, err := s.Fetch(context.Background())
			if err == nil && tc.assertErr == nil {
				testutil.AssertEq(t, "accs size", len(accs), tc.expectedSize)
				for _, acc := range accs {
					testutil.AssertNotDefault(t, "name", acc.Name)
					testutil.AssertNotDefault(t, "cpf", acc.CPF)
					testutil.AssertNotDefault(t, "balance", acc.Balance)
					testutil.AssertNotDefault(t, "created_at", acc.CreatedAt)
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
				return &testutil.AccountRepoMock{
					ExpectGetBalance: func(ctx context.Context, currentID int64) (types.Currency, error) {
						testutil.AssertEq(t, "id", id, currentID)
						return types.NewCurrency(balance), nil
					},
				}
			},
			expected:  500,
			id:        1,
			assertErr: testutil.AssertNoErr,
		},
		{
			name: "get account balance with repository error",
			repo: func(id int64, balance float64) repository.Account {
				return &testutil.AccountRepoMock{
					ExpectGetBalance: func(ctx context.Context, currentID int64) (types.Currency, error) {
						return types.NewCurrency(0), types.NewErr(types.EmptyResultErr, "no sql rows", nil)
					},
				}
			},
			expected: 0,
			id:       2,
			assertErr: func(t *testing.T, err error) {
				testutil.AssertCustomErr(t, types.EmptyResultErr, err, "no sql rows")
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo := tc.repo(tc.id, tc.expected)
			s := service.NewAccount(&txr, &repo)
			balance, err := s.GetBalance(context.Background(), tc.id)
			testutil.AssertEq(t, "balance", tc.expected, balance)
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
			expected: testutil.NewEntityAccount(1, "Sousa", "41112075020", "$2a$10$c3GzxvPAAMS9pDqB9XIYi.kT/PN7CxfRev.BsRLvAJqVcZnFiW05i", 0),
			secret:   "...",
			cpf:      "41112075020",
			repo: func(exp *entity.Account) repository.Account {
				return &testutil.AccountRepoMock{
					ExpectFindBy: func(ctx context.Context, cpf string) (*entity.Account, error) {
						testutil.AssertEq(t, "cpf", exp.CPF, cpf)
						return exp, nil
					},
				}
			},
			assertErr: testutil.AssertNoErr,
		},
		{
			name:     "login with cpf and wrong secret",
			expected: testutil.NewEntityAccount(2, "Alice", "24039310047", "123", 10),
			secret:   "123",
			cpf:      "24039310047",
			repo: func(exp *entity.Account) repository.Account {
				return &testutil.AccountRepoMock{
					ExpectFindBy: func(ctx context.Context, cpf string) (*entity.Account, error) {
						return exp, nil
					},
				}
			},
			assertErr: func(t *testing.T, err error) {
				testutil.AssertCustomErr(t, types.AuthenticationErr, err, "the provided secret doesn't match the account's secret")
			},
		},
		{
			name:   "login with repository error",
			secret: "...",
			cpf:    "72098733097",
			repo: func(exp *entity.Account) repository.Account {
				return &testutil.AccountRepoMock{
					ExpectFindBy: func(ctx context.Context, cpf string) (*entity.Account, error) {
						return exp, types.NewErr(types.InternalErr, "internal error", nil)
					},
				}
			},
			assertErr: func(t *testing.T, err error) {
				testutil.AssertCustomErr(t, types.InternalErr, err, "internal error")
			},
		},
		{
			name:   "login with nonexistent account",
			secret: "...",
			cpf:    "61632733030",
			repo: func(exp *entity.Account) repository.Account {
				return &testutil.AccountRepoMock{
					ExpectFindBy: func(ctx context.Context, cpf string) (*entity.Account, error) {
						return exp, types.NewErr(types.EmptyResultErr, "no result", nil)
					},
				}
			},
			assertErr: func(t *testing.T, err error) {
				testutil.AssertCustomErr(t, types.AuthenticationErr, err, "account with the given cpf does not exist")
			},
		},
		{
			name:   "login without cpf",
			secret: "...",
			repo: func(exp *entity.Account) repository.Account {
				return &testutil.AccountRepoMock{}
			},
			assertErr: func(t *testing.T, err error) {
				testutil.AssertCustomErr(t, types.ValidationErr, err, "field 'cpf' is required")
			},
		},
		{
			name: "login without secret",
			cpf:  "87256640005",
			repo: func(exp *entity.Account) repository.Account {
				return &testutil.AccountRepoMock{}
			},
			assertErr: func(t *testing.T, err error) {
				testutil.AssertCustomErr(t, types.ValidationErr, err, "field 'secret' is required")
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo := tc.repo(tc.expected)
			s := service.NewAccount(&txr, &repo)
			view, err := s.Login(context.Background(), tc.cpf, tc.secret)
			if err != nil {
				tc.assertErr(t, err)
			} else {
				testutil.AssertEq(t, "id", tc.expected.ID, view.ID)
				testutil.AssertEq(t, "name", tc.expected.Name, view.Name)
				testutil.AssertEq(t, "cpf", tc.expected.CPF, view.CPF)
				testutil.AssertEq(t, "balance", tc.expected.Balance.Float64(), view.Balance)
				testutil.AssertNotDefault(t, "created_at", view.CreatedAt)
			}
		})
	}
}
