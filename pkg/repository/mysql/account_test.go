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
		input []entity.Account
	}{
		{
			name: "fetch accounts with data",
			input: []entity.Account{
				testutil.NewEntityAccount(0, "Joe", "00000000000", "S001", 1),
				testutil.NewEntityAccount(0, "John", "00000000001", "S002", 2),
				testutil.NewEntityAccount(0, "Suz", "00000000002", "S003", 3),
			},
		},
		{
			name:  "fetch accounts with no result",
			input: []entity.Account{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			entities := persistTestAccountEntity(t, tc.input)

			if accs, err := repo.Fetch(context.Background()); err == nil {
				testutil.AssertEq(t, "result size", len(tc.input), len(accs))
				for _, acc := range accs {
					if expected, ok := entities[acc.ID]; ok {
						acc.ID = expected.ID
						if !reflect.DeepEqual(expected, acc) {
							t.Errorf("expected result content equal to '%v' but got '%v'", expected, acc)
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
		input []entity.Account
	}{
		{
			name: "create account with valid input",
			input: []entity.Account{
				testutil.NewEntityAccount(0, "John", "33333333331", "S300", 300),
				testutil.NewEntityAccount(0, "Jose", "33333333332", "S301", 301),
				testutil.NewEntityAccount(0, "Silva", "33333333333", "S302", 302),
				testutil.NewEntityAccount(0, "Sousa", "33333333334", "S303", 303),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(dbWipe)
			var err error
			for _, e := range tc.input {
				if id, err := repo.Create(context.Background(), e); err == nil {
					current := entity.Account{}
					row := db.QueryRow("SELECT name, cpf, secret, balance, created_at FROM account WHERE id=?", id)
					if err = row.Scan(&current.Name, &current.CPF, &current.Secret, &current.Balance, &current.CreatedAt); err == nil {
						if !reflect.DeepEqual(e, current) {
							t.Errorf("expected new account equal to '%v' but got '%v'", e, current)
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
		name   string
		input  map[int64]entity.Account
		assert func(*testing.T, error)
	}{
		{
			name: "get balance from existing account",
			input: persistTestAccountEntity(t, []entity.Account{
				testutil.NewEntityAccount(0, "Peter", "44444444440", "S400", 400),
				testutil.NewEntityAccount(0, "Tim", "44444444441", "S401", 401),
				testutil.NewEntityAccount(0, "Rom", "44444444442", "S402", 402),
			}),
			assert: func(t *testing.T, err error) {
				t.Error(err)
			},
		},
		{
			name:  "get balance from nonexisting account",
			input: map[int64]entity.Account{0: testutil.NewEntityAccount(0, "San", "44444444443", "S403", 403)},
			assert: func(t *testing.T, err error) {
				testutil.AssertCustomErr(t, types.EmptyResultErr, err, "no result getting the account balance")
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			for id, acc := range tc.input {
				if balance, err := repo.GetBalance(context.Background(), id); err == nil {
					testutil.AssertEq(t, "acc balance", acc.Balance, balance)
				} else {
					tc.assert(t, err)
				}
			}

		})
	}
}

func TestAccountRepositoryFindBy(t *testing.T) {
	repo := mysql.NewAccount(&txr)
	tt := []struct {
		name   string
		input  map[int64]entity.Account
		assert func(*testing.T, error)
	}{
		{
			name: "find by cpf with existing account",
			input: persistTestAccountEntity(t, []entity.Account{
				testutil.NewEntityAccount(0, "Maria", "55555555551", "S500", 500),
			}),
			assert: func(t *testing.T, err error) {
				t.Error(err)
			},
		},
		{
			name:  "find by cpf with nonexisting account",
			input: map[int64]entity.Account{0: testutil.NewEntityAccount(0, "Helena", "55555555552", "S501", 501)},
			assert: func(t *testing.T, err error) {
				testutil.AssertCustomErr(t, types.EmptyResultErr, err, "no result finding account by cpf")
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			for id, acc := range tc.input {
				if e, err := repo.FindBy(context.Background(), acc.CPF); err == nil {
					acc.ID = id
					if !reflect.DeepEqual(acc, e) {
						t.Errorf("expected account equal to '%v' but got '%v'", acc, e)
					}
				} else {
					tc.assert(t, err)
				}
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
		name       string
		input      map[int64]entity.Account
		assert     func(*testing.T, error)
		newBalance types.Currency
	}{
		{
			name: "update balance from existing account",
			input: persistTestAccountEntity(t, []entity.Account{
				testutil.NewEntityAccount(0, "Izzy", "66666666661", "S600", 600),
			}),
			assert: func(t *testing.T, err error) {
				t.Error(err)
			},
			newBalance: types.NewCurrency(999),
		},
		{
			name:  "update balance from nonexisting account",
			input: map[int64]entity.Account{0: testutil.NewEntityAccount(0, "Suzy", "66666666662", "S602", 602)},
			assert: func(t *testing.T, err error) {
				testutil.AssertCustomErr(t, types.NoRowAffectedErr, err, "no rows affected by the update balance stmt")
			},
			newBalance: types.NewCurrency(999),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			for id, _ := range tc.input {
				if err := repo.UpdateBalance(context.Background(), id, tc.newBalance); err == nil {

					if b, err := getCurrentAccBalance(id); err == nil {
						testutil.AssertEq(t, "new balance", tc.newBalance, b)
					} else {
						t.Error(err)
					}
				} else {
					tc.assert(t, err)
				}
			}

		})
	}
}

func TestAccountRepositoryExists(t *testing.T) {
	repo := mysql.NewAccount(&txr)
	tt := []struct {
		name     string
		input    map[int64]entity.Account
		expected bool
	}{
		{
			name: "existing account",
			input: persistTestAccountEntity(t, []entity.Account{
				testutil.NewEntityAccount(0, "William", "77777777771", "S701", 701),
			}),
			expected: true,
		},
		{
			name:  "nonexisting account",
			input: map[int64]entity.Account{0: testutil.NewEntityAccount(0, "James", "77777777771", "S702", 702)},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			for id, _ := range tc.input {
				if exists, err := repo.Exists(context.Background(), id); err == nil {
					testutil.AssertEq(t, "account existence", tc.expected, exists)
				} else {
					t.Error(err)
				}
			}

		})
	}
}
