package validation_test

import (
	"context"
	"testing"

	"github.com/rafael-sousa/stn-accounts/pkg/model/dto"
	"github.com/rafael-sousa/stn-accounts/pkg/model/entity"
	"github.com/rafael-sousa/stn-accounts/pkg/model/types"
)

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

func newAccountCreation(n, cpf, s string, b float64) *dto.AccountCreation {
	return &dto.AccountCreation{
		Name:    n,
		CPF:     cpf,
		Secret:  s,
		Balance: b,
	}
}

func assertCustomErr(t *testing.T, c types.ErrCode, err error, msg string) {
	if customErr, ok := err.(*types.Err); ok {
		assertEq(t, "err code", c, customErr.Code)
		assertEq(t, "err msg", msg, customErr.Msg)
	} else {
		t.Errorf("expected err equal to a types.Err but got %v", err)
	}
}

func assertEq(t *testing.T, n string, expected interface{}, current interface{}) {
	if expected != current {
		t.Errorf("expected %s equal to '%v' but got '%v'", n, expected, current)
	}
}
