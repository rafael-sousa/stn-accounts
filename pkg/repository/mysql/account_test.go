package mysql_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/rafael-sousa/stn-accounts/pkg/model/entity"
	"github.com/rafael-sousa/stn-accounts/pkg/model/types"
	"github.com/rafael-sousa/stn-accounts/pkg/repository/mysql"
)

func persistTestAccountEntity(t *testing.T, input []*entity.Account) map[int64]*entity.Account {
	entities := make(map[int64]*entity.Account, 0)
	if len(input) == 0 {
		return entities
	}
	t.Cleanup(dbWipe)
	stmt, err := db.Prepare("INSERT INTO account(name,cpf,secret,balance,created_at) VALUES (?,?,?,?,?)")
	logFatal(err, "unable to prepare account insert stmt")
	defer stmt.Close()
	for _, e := range input {
		result, _ := stmt.Exec(e.Name, e.CPF, e.Secret, e.Balance, e.CreatedAt)
		logFatal(err, "unable to exec account insert stmt")
		id, _ := result.LastInsertId()
		logFatal(err, "unable to retrieve inserted account id")
		e.ID = id
		entities[id] = e
	}
	return entities
}

func TestAccountRepositoryFetch(t *testing.T) {
	repo := mysql.NewAccount(&txr)
	tt := []struct {
		name  string
		input []*entity.Account
	}{
		{
			name: "fetch accounts with data",
			input: []*entity.Account{
				newAccount(0, "Joe", "00000000000", "S001", 1),
				newAccount(0, "John", "00000000001", "S002", 2),
				newAccount(0, "Suz", "00000000002", "S003", 3),
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
				if len(accs) != len(tc.input) {
					t.Errorf("expected result size equal to '%d' but got '%d'", len(tc.input), len(accs))
				}
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
				newAccount(0, "John", "33333333331", "S300", 300),
				newAccount(0, "Jose", "33333333332", "S301", 301),
				newAccount(0, "Silva", "33333333333", "S302", 302),
				newAccount(0, "Sousa", "33333333334", "S303", 303),
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
			input: newAccount(0, "Peter", "44444444441", "S400", 400),
			prepare: func(t *testing.T, e *entity.Account) {
				persistTestAccountEntity(t, []*entity.Account{e})
			},
			assert: func(t *testing.T, err error) {
				t.Error(err)
			},
		},
		{
			name:    "get balance from nonexisting account",
			input:   newAccount(0, "Bob", "44444444442", "S401", 401),
			prepare: func(t *testing.T, e *entity.Account) {},
			assert: func(t *testing.T, err error) {
				if customErr, ok := err.(*types.Err); ok {
					msg := "no result getting the account balance"
					if customErr.Msg != msg {
						t.Errorf("expected error message equal to '%s' but got '%s'", msg, customErr.Msg)
					}
					if customErr.Code != types.EmptyResultErr {
						t.Errorf("expected error code equal to '%v' but got '%v'", types.EmptyResultErr, customErr.Code)
					}
				} else {
					t.Error(err)
				}
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
			input: newAccount(0, "Maria", "55555555551", "S500", 500),
			prepare: func(t *testing.T, e *entity.Account) {
				persistTestAccountEntity(t, []*entity.Account{e})
			},
			assert: func(t *testing.T, err error) {
				t.Error(err)
			},
		},
		{
			name:    "find by cpf with nonexisting account",
			input:   newAccount(0, "Helena", "55555555552", "S501", 501),
			prepare: func(t *testing.T, e *entity.Account) {},
			assert: func(t *testing.T, err error) {
				if customErr, ok := err.(*types.Err); ok {
					msg := "no result finding account by cpf"
					if customErr.Msg != msg {
						t.Errorf("expected error message equal to '%s' but got '%s'", msg, customErr.Msg)
					}
					if customErr.Code != types.EmptyResultErr {
						t.Errorf("expected error code equal to '%v' but got '%v'", types.EmptyResultErr, customErr.Code)
					}
				} else {
					t.Error(err)
				}
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
			input: newAccount(0, "Izzy", "66666666661", "S600", 600),
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
			input:   newAccount(0, "Suzy", "66666666662", "S602", 602),
			prepare: func(t *testing.T, e *entity.Account) {},
			assert: func(t *testing.T, err error) {
				if customErr, ok := err.(*types.Err); ok {
					msg := "no rows affected by the update balance stmt"
					if customErr.Msg != msg {
						t.Errorf("expected error message equal to '%s' but got '%s'", msg, customErr.Msg)
					}
					if customErr.Code != types.NoRowAffectedErr {
						t.Errorf("expected error code equal to '%v' but got '%v'", types.NoRowAffectedErr, customErr.Code)
					}
				} else {
					t.Error(err)
				}
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
