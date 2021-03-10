package testutil

import (
	"context"

	"github.com/rafael-sousa/stn-accounts/pkg/model/entity"
	"github.com/rafael-sousa/stn-accounts/pkg/model/types"
	"github.com/rafael-sousa/stn-accounts/pkg/repository"
)

// TransactionerMock mocks the repository.Transactioner interface
type TransactionerMock struct{}

// WithTx mocks the transactional context and call the fuction stored at fn
func (manager *TransactionerMock) WithTx(ctx context.Context, fn func(context.Context) error) (err error) {
	return fn(ctx)
}

// GetConn represents the default mock of repository.Transactioner#GetConn
func (manager *TransactionerMock) GetConn(ctx context.Context) repository.Connection {
	return nil
}

// AccountRepoMock mock structure for repository.Account interface
type AccountRepoMock struct {
	ExpectFetch         func(context.Context) ([]*entity.Account, error)
	ExpectCreate        func(context.Context, *entity.Account) (*entity.Account, error)
	ExpectGetBalance    func(context.Context, int64) (types.Currency, error)
	ExpectFindBy        func(context.Context, string) (*entity.Account, error)
	ExpectUpdateBalance func(context.Context, int64, types.Currency) error
	ExpectExists        func(context.Context, int64) (bool, error)
}

// Fetch mocks the functionality of repository.Account#Fetch
func (r *AccountRepoMock) Fetch(ctx context.Context) ([]*entity.Account, error) {
	return r.ExpectFetch(ctx)
}

// Create mocks the functionality of repository.Account#Create
func (r *AccountRepoMock) Create(ctx context.Context, e *entity.Account) (*entity.Account, error) {
	return r.ExpectCreate(ctx, e)
}

// GetBalance mocks the functionality of repository.Account#GetBalance
func (r *AccountRepoMock) GetBalance(ctx context.Context, id int64) (types.Currency, error) {
	return r.ExpectGetBalance(ctx, id)
}

// FindBy mocks the functionality of repository.Account#FindBy
func (r *AccountRepoMock) FindBy(ctx context.Context, cpf string) (*entity.Account, error) {
	return r.ExpectFindBy(ctx, cpf)
}

// UpdateBalance mocks the functionality of repository.Account#UpdateBalance
func (r *AccountRepoMock) UpdateBalance(ctx context.Context, id int64, b types.Currency) error {
	return r.ExpectUpdateBalance(ctx, id, b)
}

// Exists mocks the functionality of repository.Account#Exists
func (r *AccountRepoMock) Exists(ctx context.Context, id int64) (bool, error) {
	return r.ExpectExists(ctx, id)
}

// TransferRepoMock mocks the repository.Transfer interface
type TransferRepoMock struct {
	ExpectFetch  func(ctx context.Context, id int64) ([]*entity.Transfer, error)
	ExpectCreate func(ctx context.Context, e *entity.Transfer) (*entity.Transfer, error)
}

// Fetch mocks the functionality of repository.Transfer#Fetch
func (r *TransferRepoMock) Fetch(ctx context.Context, id int64) ([]*entity.Transfer, error) {
	return r.ExpectFetch(ctx, id)

}

// Create mocks the functionality of repository.Transfer#Create
func (r *TransferRepoMock) Create(ctx context.Context, e *entity.Transfer) (*entity.Transfer, error) {
	return r.ExpectCreate(ctx, e)
}
