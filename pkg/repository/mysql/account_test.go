package mysql_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/rafael-sousa/stn-accounts/pkg/model/entity"
	"github.com/rafael-sousa/stn-accounts/pkg/model/types"
	"github.com/rafael-sousa/stn-accounts/pkg/repository/mysql"
)

func createTestAccountEntity(t *testing.T) int64 {
	result, err := db.Exec("INSERT INTO account(name,cpf,secret,balance,created_at) VALUES ('John','99999999999','pw',500,'2100-12-31')")
	logFatal(err, "unable to prepare testcase")
	t.Cleanup(dbWipe)
	id, err := result.LastInsertId()
	logFatal(err, "unable to retrieve inserted id")
	return id
}

func TestAccountRepositoryFetch(t *testing.T) {
	repo := mysql.NewAccount(&txr)
	tt := []struct {
		name         string
		expectedSize int64
		prepare      func(*testing.T) int64
	}{
		{
			name:         "fetch accounts with result",
			expectedSize: 1,
			prepare:      createTestAccountEntity,
		},
		{
			name:         "fetch accounts with no results",
			expectedSize: 0,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.prepare != nil {
				tc.prepare(t)
			}

			if accs, err := repo.Fetch(context.Background()); err == nil {
				if len(accs) != int(tc.expectedSize) {
					t.Errorf("expected return size equal to '%d' but got '%d'", tc.expectedSize, len(accs))
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
		name string
		e    *entity.Account
	}{
		{
			name: "create account successfully",
			e: &entity.Account{
				Name:      "John",
				CPF:       "00000000000",
				Secret:    "pw",
				Balance:   100,
				CreatedAt: time.Now(),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(dbWipe)
			if e, err := repo.Create(context.Background(), tc.e); err == nil {
				if !reflect.DeepEqual(tc.e, e) {
					t.Errorf("expected new account equal to '%v' but got '%v'", tc.e, e)
				}
			} else {
				t.Error(err)
			}
		})
	}
}

func TestAccountRepositoryGetBalance(t *testing.T) {
	repo := mysql.NewAccount(&txr)
	tt := []struct {
		name     string
		expected types.Currency
		prepare  func(*testing.T) int64
	}{
		{
			name:     "get account balance",
			expected: types.NewCurrency(5),
			prepare:  createTestAccountEntity,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			id := tc.prepare(t)
			if balance, err := repo.GetBalance(context.Background(), id); err == nil {
				if balance != tc.expected {
					t.Errorf("expected balance equal to '%d' but got '%d'", tc.expected, balance)
				}
			} else {
				t.Error(err)
			}

		})
	}
}

func TestAccountRepositoryFindBy(t *testing.T) {
	repo := mysql.NewAccount(&txr)
	tt := []struct {
		name     string
		expected string
		prepare  func(*testing.T) int64
	}{
		{
			name:     "find account by cpf",
			expected: "99999999999",
			prepare:  createTestAccountEntity,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare(t)

			if acc, err := repo.FindBy(context.Background(), tc.expected); err == nil {
				if acc.CPF != tc.expected {
					t.Errorf("expected cpf equal to '%s' but got '%s'", tc.expected, acc.CPF)
				}
			} else {
				t.Error(err)
			}

		})
	}
}

func TestAccountRepositoryUpdateBalance(t *testing.T) {
	repo := mysql.NewAccount(&txr)

	getCurrentAccBalance := func(id int64) types.Currency {
		var balance types.Currency
		err := db.QueryRow("SELECT balance FROM account WHERE id=?", id).Scan(&balance)
		logFatal(err, "unable to query account balance")
		return balance
	}
	tt := []struct {
		name     string
		expected types.Currency
		prepare  func(*testing.T) int64
	}{
		{
			name:     "update account balance",
			expected: types.NewCurrency(99),
			prepare:  createTestAccountEntity,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			id := tc.prepare(t)

			if err := repo.UpdateBalance(context.Background(), id, tc.expected); err == nil {
				if balance := getCurrentAccBalance(id); balance != tc.expected {
					t.Errorf("expected updated balance '%d' but got '%d'", tc.expected, balance)
				}
			} else {
				t.Error(err)
			}

		})
	}
}
