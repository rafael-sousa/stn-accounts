package routing_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/rafael-sousa/stn-accounts/pkg/controller/rest/jwt"
	"github.com/rafael-sousa/stn-accounts/pkg/model/dto"
	"github.com/rafael-sousa/stn-accounts/pkg/model/env"
)

type accountServMock struct {
	fetch      func(context.Context) ([]*dto.AccountView, error)
	getBalance func(context.Context, int64) (float64, error)
	create     func(context.Context, *dto.AccountCreation) (*dto.AccountView, error)
	login      func(context.Context, string, string) (*dto.AccountView, error)
}

func (s *accountServMock) Fetch(ctx context.Context) ([]*dto.AccountView, error) {
	return s.fetch(ctx)
}
func (s *accountServMock) GetBalance(ctx context.Context, id int64) (float64, error) {
	return s.getBalance(ctx, id)
}
func (s *accountServMock) Create(ctx context.Context, d *dto.AccountCreation) (*dto.AccountView, error) {
	return s.create(ctx, d)
}
func (s *accountServMock) Login(ctx context.Context, cpf string, secret string) (*dto.AccountView, error) {
	return s.login(ctx, cpf, secret)
}

type transferServMock struct {
	fetch  func(context.Context, int64) ([]*dto.TransferView, error)
	create func(context.Context, int64, *dto.TransferCreation) (*dto.TransferView, error)
}

func (s *transferServMock) Fetch(ctx context.Context, id int64) ([]*dto.TransferView, error) {
	return s.fetch(ctx, id)
}
func (s *transferServMock) Create(ctx context.Context, origin int64, d *dto.TransferCreation) (*dto.TransferView, error) {
	return s.create(ctx, origin, d)
}

func newAccountView(id int64, name, cpf string, balance float64, createdAt time.Time) *dto.AccountView {
	return &dto.AccountView{
		ID:        id,
		Name:      name,
		CPF:       cpf,
		Balance:   balance,
		CreatedAt: createdAt,
	}
}

func newAccountCreation(name, cpf string, balance float64) *dto.AccountCreation {
	return &dto.AccountCreation{
		Name:    name,
		CPF:     cpf,
		Balance: balance,
	}
}

func newTransferCreation(dest int64, amt float64) *dto.TransferCreation {
	return &dto.TransferCreation{
		Destination: dest,
		Amount:      amt,
	}
}

func newTransferView(id, destination int64, amount float64) *dto.TransferView {
	return &dto.TransferView{
		ID:          id,
		Destination: destination,
		Amount:      amount,
		CreatedAt:   time.Now(),
	}
}

var jwtHandler *jwt.Handler

func TestMain(m *testing.M) {
	cfg := env.RestConfig{
		Secret:          []byte("secret"),
		TokenExpTimeout: 30,
	}
	jwtHandler = jwt.NewHandler(&cfg)
	os.Exit(m.Run())
}

func assertEq(t *testing.T, n string, expected interface{}, current interface{}) {
	if expected != current {
		t.Errorf("expected %s equal to '%v' but got '%v'", n, expected, current)
	}
}
