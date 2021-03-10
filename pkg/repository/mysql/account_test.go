package mysql_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/rafael-sousa/stn-accounts/pkg/model/entity"
	"github.com/rafael-sousa/stn-accounts/pkg/model/types"
	"github.com/rafael-sousa/stn-accounts/pkg/repository/mysql"
	"github.com/rafael-sousa/stn-accounts/pkg/testutil"
)

func TestAccountRepositoryFetch(t *testing.T) {
	repo := mysql.NewAccount(&txr)
	tt := []struct {
		name  string
		input []*entity.Account
	}{
		{
			name: "fetch accounts with data",
			input: []*entity.Account{
				testutil.NewAccount(0, "Joe", "00000000000", "S001", 1),
				testutil.NewAccount(0, "John", "00000000001", "S002", 2),
				testutil.NewAccount(0, "Suz", "00000000002", "S003", 3),
			},
		},
		{
			name:  "fetch accounts with no result",
			input: []*entity.Account{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			entities := persistTestAccountEntity(t, tc.input)

			if accs, err := repo.Fetch(context.Background()); err == nil {
				testutil.AssertEq(t, "result size", len(tc.input), len(accs))
				for _, acc := range accs {
					if expected, ok := entities[acc.ID]; ok {
						if !reflect.DeepEqual(*expected, *acc) {
							t.Errorf("expected result content equal to '%v' but got '%v'", *expected, *acc)
						}
					} else {
						t.Errorf("unexpected account id '%d'", acc.ID)
					}
				}
			} else {
				t.Error(err)
			}
		})
	}
}

func TestAccountRepositoryCreate(t *testing.T) {
	repo := mysql.NewAccount(&txr)
	tt := []struct {
		name  string
		input []*entity.Account
	}{
		{
			name: "create account with valid input",
			input: []*entity.Account{
				testutil.NewAccount(0, "John", "33333333331", "S300", 300),
				testutil.NewAccount(0, "Jose", "33333333332", "S301", 301),
				testutil.NewAccount(0, "Silva", "33333333333", "S302", 302),
				testutil.NewAccount(0, "Sousa", "33333333334", "S303", 303),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(dbWipe)
			var err error
			for _, e := range tc.input {
				if _, err = repo.Create(context.Background(), e); err == nil {
					current := entity.Account{}
					row := db.QueryRow("SELECT id, name, cpf, secret, balance, created_at FROM account WHERE id=?", e.ID)
					if err = row.Scan(&current.ID, &current.Name, &current.CPF, &current.Secret, &current.Balance, &current.CreatedAt); err == nil {
						if !reflect.DeepEqual(*e, current) {
							t.Errorf("expected new account equal to '%v' but got '%v'", *e, current)
						}
					}
				}
				if err != nil {
					t.Error(err)
				}
			}

		})
	}
}

func TestAccountRepositoryGetBalance(t *testing.T) {
	repo := mysql.NewAccount(&txr)
	tt := []struct {
		name    string
		input   *entity.Account
		prepare func(*testing.T, *entity.Account)
		assert  func(*testing.T, error)
	}{
		{
			name:  "get balance from existing account",
			input: testutil.NewAccount(0, "Peter", "44444444441", "S400", 400),
			prepare: func(t *testing.T, e *entity.Account) {
				persistTestAccountEntity(t, []*entity.Account{e})
			},
			assert: func(t *testing.T, err error) {
				t.Error(err)
			},
		},
		{
			name:    "get balance from nonexisting account",
			input:   testutil.NewAccount(0, "Bob", "44444444442", "S401", 401),
			prepare: func(t *testing.T, e *entity.Account) {},
			assert: func(t *testing.T, err error) {
				testutil.AssertCustomErr(t, types.EmptyResultErr, err, "no result getting the account balance")
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare(t, tc.input)
			if b, err := repo.GetBalance(context.Background(), tc.input.ID); err == nil {
				if b != tc.input.Balance {
					t.Errorf("expected balance equal to '%d' but got '%d'", tc.input.Balance, b)
				}
			} else {
				tc.assert(t, err)
			}
		})
	}
}

func TestAccountRepositoryFindBy(t *testing.T) {
	repo := mysql.NewAccount(&txr)
	tt := []struct {
		name    string
		input   *entity.Account
		prepare func(*testing.T, *entity.Account)
		assert  func(*testing.T, error)
	}{
		{
			name:  "find by cpf with existing account",
			input: testutil.NewAccount(0, "Maria", "55555555551", "S500", 500),
			prepare: func(t *testing.T, e *entity.Account) {
				persistTestAccountEntity(t, []*entity.Account{e})
			},
			assert: func(t *testing.T, err error) {
				t.Error(err)
			},
		},
		{
			name:    "find by cpf with nonexisting account",
			input:   testutil.NewAccount(0, "Helena", "55555555552", "S501", 501),
			prepare: func(t *testing.T, e *entity.Account) {},
			assert: func(t *testing.T, err error) {
				testutil.AssertCustomErr(t, types.EmptyResultErr, err, "no result finding account by cpf")
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare(t, tc.input)
			if e, err := repo.FindBy(context.Background(), tc.input.CPF); err == nil {

				if !reflect.DeepEqual(*tc.input, *e) {
					t.Errorf("expected account equal to '%v' but got '%v'", *tc.input, *e)
				}
			} else {
				tc.assert(t, err)
			}
		})
	}
}

func TestAccountRepositoryUpdateBalance(t *testing.T) {
	getCurrentAccBalance := func(id int64) (types.Currency, error) {
		var balance types.Currency
		err := db.QueryRow("SELECT balance FROM account WHERE id=?", id).Scan(&balance)
		return balance, err
	}

	repo := mysql.NewAccount(&txr)
	tt := []struct {
		name            string
		input           *entity.Account
		prepare         func(*testing.T, *entity.Account)
		assert          func(*testing.T, error)
		expectedBalance types.Currency
	}{
		{
			name:  "update balance from existing account",
			input: testutil.NewAccount(0, "Izzy", "66666666661", "S600", 600),
			prepare: func(t *testing.T, e *entity.Account) {
				persistTestAccountEntity(t, []*entity.Account{e})
			},
			assert: func(t *testing.T, err error) {
				t.Error(err)
			},
			expectedBalance: types.NewCurrency(999),
		},
		{
			name:    "update balance from nonexisting account",
			input:   testutil.NewAccount(0, "Suzy", "66666666662", "S602", 602),
			prepare: func(t *testing.T, e *entity.Account) {},
			assert: func(t *testing.T, err error) {
				testutil.AssertCustomErr(t, types.NoRowAffectedErr, err, "no rows affected by the update balance stmt")
			},
			expectedBalance: types.NewCurrency(999),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare(t, tc.input)
			if err := repo.UpdateBalance(context.Background(), tc.input.ID, tc.expectedBalance); err == nil {

				if b, err := getCurrentAccBalance(tc.input.ID); err == nil {
					if b != tc.expectedBalance {
						t.Errorf("expected updated account balance equal to '%v' but got '%v'", tc.expectedBalance, b)
					}
				} else {
					t.Error(err)
				}
			} else {
				tc.assert(t, err)
			}
		})
	}
}

func TestAccountRepositoryExists(t *testing.T) {
	repo := mysql.NewAccount(&txr)
	tt := []struct {
		name     string
		input    *entity.Account
		prepare  func(*testing.T, *entity.Account)
		expected bool
	}{
		{
			name:  "existing account",
			input: testutil.NewAccount(0, "William", "77777777771", "S701", 701),
			prepare: func(t *testing.T, e *entity.Account) {
				persistTestAccountEntity(t, []*entity.Account{e})
			},
			expected: true,
		},
		{
			name:    "nonexisting account",
			input:   testutil.NewAccount(0, "James", "77777777771", "S702", 702),
			prepare: func(t *testing.T, e *entity.Account) {},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare(t, tc.input)
			if exists, err := repo.Exists(context.Background(), tc.input.ID); err == nil {
				testutil.AssertEq(t, "account existence", tc.expected, exists)
			} else {
				t.Error(err)
			}
		})
	}
}
