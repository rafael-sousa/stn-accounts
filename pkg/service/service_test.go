package service_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/rafael-sousa/stn-accounts/pkg/model/dto"
	"github.com/rafael-sousa/stn-accounts/pkg/model/entity"
	"github.com/rafael-sousa/stn-accounts/pkg/model/types"
	"github.com/rafael-sousa/stn-accounts/pkg/repository"
)

type transactionerMock struct{}

func (manager *transactionerMock) WithTx(ctx context.Context, f func(context.Context) error) (err error) {
	return f(ctx)
}

func (manager *transactionerMock) GetConn(ctx context.Context) repository.Connection {
	return nil
}

type accountRepoMock struct {
	fetch         func(context.Context) ([]*entity.Account, error)
	create        func(context.Context, *entity.Account) (*entity.Account, error)
	getBalance    func(context.Context, int64) (types.Currency, error)
	findBy        func(context.Context, string) (*entity.Account, error)
	updateBalance func(context.Context, int64, types.Currency) error
	exists        func(context.Context, int64) (bool, error)
}

func (r *accountRepoMock) Fetch(ctx context.Context) ([]*entity.Account, error) {
	return r.fetch(ctx)
}
func (r *accountRepoMock) Create(ctx context.Context, e *entity.Account) (*entity.Account, error) {
	return r.create(ctx, e)
}
func (r *accountRepoMock) GetBalance(ctx context.Context, id int64) (types.Currency, error) {
	return r.getBalance(ctx, id)
}
func (r *accountRepoMock) FindBy(ctx context.Context, cpf string) (*entity.Account, error) {
	return r.findBy(ctx, cpf)
}
func (r *accountRepoMock) UpdateBalance(ctx context.Context, id int64, b types.Currency) error {
	return r.updateBalance(ctx, id, b)
}

func (r *accountRepoMock) Exists(ctx context.Context, id int64) (bool, error) {
	return r.exists(ctx, id)
}

type transferRepoMock struct {
	fetch  func(ctx context.Context, id int64) ([]*entity.Transfer, error)
	create func(ctx context.Context, e *entity.Transfer) (*entity.Transfer, error)
}

func (r *transferRepoMock) Fetch(ctx context.Context, id int64) ([]*entity.Transfer, error) {
	return r.fetch(ctx, id)

}
func (r *transferRepoMock) Create(ctx context.Context, e *entity.Transfer) (*entity.Transfer, error) {
	return r.create(ctx, e)
}

var txr repository.Transactioner

func assertEq(t *testing.T, n string, expected interface{}, current interface{}) {
	if expected != current {
		t.Errorf("expected %s equal to '%v' but got '%v'", n, expected, current)
	}
}

func assertNotDefault(t *testing.T, n string, current interface{}) {
	switch current.(type) {
	case time.Time:
		if current.(time.Time).IsZero() {
			t.Errorf("expected %s not empty", n)
		}
	case string:
		if len(current.(string)) == 0 {
			t.Errorf("expected %s not empty", n)
		}
	case int, int32, int64, float32, float64:
		if current == 0 {
			t.Errorf("expected %s not zero", n)
		}
	case nil:
		t.Errorf("expected %s not nil", n)
	}
}

func assertCustomErr(t *testing.T, c types.ErrCode, err error, msg string) {
	if customErr, ok := err.(*types.Err); ok {
		assertEq(t, "err code", c, customErr.Code)
		assertEq(t, "err msg", customErr.Msg, msg)
	} else {
		t.Errorf("expected err equal to a types.Err but got %v", err)
	}
}

func newAccountCreation(n, cpf, s string, b float64) *dto.AccountCreation {
	return &dto.AccountCreation{
		Name:    n,
		CPF:     cpf,
		Secret:  s,
		Balance: b,
	}
}

func newAccount(id int64, n, cpf, s string, b float64) *entity.Account {
	return &entity.Account{
		ID:        id,
		Name:      n,
		CPF:       cpf,
		Secret:    s,
		Balance:   types.NewCurrency(b),
		CreatedAt: time.Now(),
	}
}

func newTransfer(id, origin, destination int64, b float64) *entity.Transfer {
	return &entity.Transfer{
		ID:          id,
		Origin:      origin,
		Destination: destination,
		Amount:      types.NewCurrency(b),
		CreatedAt:   time.Now(),
	}
}

func TestMain(m *testing.M) {
	txr = &transactionerMock{}
	os.Exit(m.Run())
}
