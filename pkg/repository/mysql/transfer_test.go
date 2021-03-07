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

func TestTransferRepositoryFetch(t *testing.T) {
	repo := mysql.NewTransfer(&txr)
	tt := []struct {
		name         string
		expectedSize int
		prepare      func(*testing.T) (int64, int64, int64)
	}{
		{
			name:         "fetch transfers with result",
			expectedSize: 1,
			prepare: func(t *testing.T) (int64, int64, int64) {
				t.Cleanup(dbWipe)

				result, err := db.Exec("INSERT INTO account(name,cpf,secret,balance,created_at) VALUES ('John','99999999999','pw',100,'2100-12-31')")
				logFatal(err, "unable to prepare testcase")

				origin, err := result.LastInsertId()
				logFatal(err, "unable to retrieve inserted id")

				result, err = db.Exec("INSERT INTO account(name,cpf,secret,balance,created_at) VALUES ('Doe','88888888888','pw',500,'2100-12-31')")
				logFatal(err, "unable to prepare testcase")

				destination, err := result.LastInsertId()
				logFatal(err, "unable to retrieve inserted id")

				stmt, err := db.Prepare("INSERT INTO transfer(account_origin_id, account_destination_id, amount, created_at) VALUES (?,?,?,?)")
				logFatal(err, "unable to prepare testcase")
				result, err = stmt.Exec(origin, destination, 500, time.Now())

				id, err := result.LastInsertId()
				logFatal(err, "unable to retrieve inserted id")
				return origin, destination, id
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			origin, _, _ := tc.prepare(t)
			if transfers, err := repo.Fetch(context.Background(), origin); err == nil {
				if len(transfers) != tc.expectedSize {
					t.Errorf("expected return size equal to '%d' but got '%d'", tc.expectedSize, len(transfers))
				}
			} else {
				t.Error(err)
			}

		})
	}
}

func TestTransferRepositoryCreate(t *testing.T) {
	repo := mysql.NewTransfer(&txr)
	tt := []struct {
		name         string
		expectedSize int64
		prepare      func(*testing.T) *entity.Transfer
	}{
		{
			name:         "create transfer with result",
			expectedSize: 1,
			prepare: func(t *testing.T) *entity.Transfer {
				t.Cleanup(dbWipe)

				result, err := db.Exec("INSERT INTO account(name,cpf,secret,balance,created_at) VALUES ('John','99999999999','pw',100,'2100-12-31')")
				logFatal(err, "unable to prepare testcase")

				origin, err := result.LastInsertId()
				logFatal(err, "unable to retrieve inserted id")

				result, err = db.Exec("INSERT INTO account(name,cpf,secret,balance,created_at) VALUES ('Doe','88888888888','pw',0,'2100-12-31')")
				logFatal(err, "unable to prepare testcase")

				destination, err := result.LastInsertId()
				logFatal(err, "unable to retrieve inserted id")

				return &entity.Transfer{
					Origin:      origin,
					Destination: destination,
					Amount:      types.NewCurrency(5),
				}
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			e := tc.prepare(t)
			if transfer, err := repo.Create(context.Background(), e); err == nil {
				if !reflect.DeepEqual(e, transfer) {
					t.Errorf("expected new transfer equal to '%v' but got '%v'", e, transfer)
				}
			} else {
				t.Error(err)
			}

		})
	}
}
