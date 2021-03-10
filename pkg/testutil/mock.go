package testutil

import (
	"context"

	"github.com/rafael-sousa/stn-accounts/pkg/model/dto"
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

// AccountServMock mocks the service.Account interface
type AccountServMock struct {
	ExpectFetch      func(context.Context) ([]*dto.AccountView, error)
	ExpectGetBalance func(context.Context, int64) (float64, error)
	ExpectCreate     func(context.Context, *dto.AccountCreation) (*dto.AccountView, error)
	ExpectLogin      func(context.Context, string, string) (*dto.AccountView, error)
}

// Fetch mocks the functionality of service.Account#Fetch
func (s *AccountServMock) Fetch(ctx context.Context) ([]*dto.AccountView, error) {
	return s.ExpectFetch(ctx)
}

// GetBalance mocks the functionality of service.Account#GetBalance
func (s *AccountServMock) GetBalance(ctx context.Context, id int64) (float64, error) {
	return s.ExpectGetBalance(ctx, id)
}

// Create mocks the functionality of service.Account#Create
func (s *AccountServMock) Create(ctx context.Context, d *dto.AccountCreation) (*dto.AccountView, error) {
	return s.ExpectCreate(ctx, d)
}

// Login mocks the functionality of service.Account#Login
func (s *AccountServMock) Login(ctx context.Context, cpf string, secret string) (*dto.AccountView, error) {
	return s.ExpectLogin(ctx, cpf, secret)
}

// TransferServMock mocks the service.Transfer interface
type TransferServMock struct {
	ExpectFetch  func(context.Context, int64) ([]*dto.TransferView, error)
	ExpectCreate func(context.Context, int64, *dto.TransferCreation) (*dto.TransferView, error)
}

// Fetch mocks the functionality of service.Transfer#Fetch
func (s *TransferServMock) Fetch(ctx context.Context, id int64) ([]*dto.TransferView, error) {
	return s.ExpectFetch(ctx, id)
}

// Create mocks the functionality of service.Transfer#Create
func (s *TransferServMock) Create(ctx context.Context, origin int64, d *dto.TransferCreation) (*dto.TransferView, error) {
	return s.ExpectCreate(ctx, origin, d)
}
