// Package mysql contains artefacts that implements the repository interfaces for the mysql db
package mysql

import (
	"context"
	"database/sql"

	"github.com/rafael-sousa/stn-accounts/pkg/model/entity"
	"github.com/rafael-sousa/stn-accounts/pkg/model/types"
	"github.com/rafael-sousa/stn-accounts/pkg/repository"
)

type account struct {
	txr *repository.Transactioner
}

var _ repository.Account = (*account)(nil)

// NewAccount creates a value that satisfies the repository.Account interface
func NewAccount(txr *repository.Transactioner) repository.Account {
	return &account{txr: txr}
}

func (r *account) Fetch(ctx context.Context) ([]entity.Account, error) {
	rows, err := (*r.txr).GetConn(ctx).QueryContext(ctx, "SELECT id, name, cpf, secret, balance, created_at FROM account")
	if err != nil {
		return nil, types.NewErr(types.SelectStmtErr, "fetching accounts", err)
	}
	defer rows.Close()
	accs := make([]entity.Account, 0)
	for rows.Next() {
		acc := entity.Account{}
		if err = rows.Scan(&acc.ID, &acc.Name, &acc.CPF, &acc.Secret, &acc.Balance, &acc.CreatedAt); err != nil {
			return nil, types.NewErr(types.SelectStmtErr, "scanning account row", err)
		}
		accs = append(accs, acc)
	}

	if err = rows.Err(); err != nil {
		return nil, types.NewErr(types.SelectStmtErr, "iterating over the account rows", err)
	}
	return accs, nil
}

func (r *account) Create(ctx context.Context, e entity.Account) (insertedID int64, err error) {
	stmt, err := (*r.txr).GetConn(ctx).PrepareContext(ctx, "INSERT INTO account(name, cpf, secret, balance, created_at) VALUES (?,?,?,?,?)")
	if err != nil {
		return insertedID, types.NewErr(types.InsertStmtErr, "preparing account insert stmt", err)
	}
	defer stmt.Close()
	result, err := stmt.ExecContext(ctx, e.Name, e.CPF, e.Secret, e.Balance, e.CreatedAt)
	if err != nil {
		return insertedID, types.NewErr(types.InsertStmtErr, "exec account insert stmt", err)
	}
	if insertedID, err = result.LastInsertId(); err != nil {
		return insertedID, types.NewErr(types.InsertStmtErr, "getting the inserted account id", err)
	}
	return insertedID, nil
}
func (r *account) GetBalance(ctx context.Context, id int64) (types.Currency, error) {
	var balance types.Currency
	err := (*r.txr).GetConn(ctx).QueryRowContext(ctx, "SELECT balance FROM account WHERE id=?", id).Scan(&balance)
	if err == sql.ErrNoRows {
		return 0, types.NewErr(types.EmptyResultErr, "no result getting the account balance", nil)
	}
	if err != nil {
		return 0, types.NewErr(types.SelectStmtErr, "scanning the account balance row", err)
	}
	return balance, nil
}

func (r *account) FindBy(ctx context.Context, cpf string) (acc entity.Account, err error) {
	q := "SELECT id, name, cpf, secret, balance, created_at FROM account WHERE cpf=?"
	err = (*r.txr).GetConn(ctx).QueryRowContext(ctx, q, cpf).Scan(&acc.ID, &acc.Name, &acc.CPF, &acc.Secret, &acc.Balance, &acc.CreatedAt)
	if err == sql.ErrNoRows {
		return acc, types.NewErr(types.EmptyResultErr, "no result finding account by cpf", err)
	}
	if err != nil {
		return acc, types.NewErr(types.SelectStmtErr, "finding account by cpf", err)
	}
	return acc, nil
}

func (r *account) UpdateBalance(ctx context.Context, id int64, balance types.Currency) error {
	stmt, err := (*r.txr).GetConn(ctx).PrepareContext(ctx, "UPDATE account SET balance=? WHERE id=?")
	if err != nil {
		return types.NewErr(types.UpdateStmtErr, "preparing update account balance stmt", err)
	}
	defer stmt.Close()
	result, err := stmt.ExecContext(ctx, balance, id)
	if err != nil {
		return types.NewErr(types.UpdateStmtErr, "exec the update account balance stmt", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return types.NewErr(types.UpdateStmtErr, "getting the affected rows number", err)
	}
	if rowsAffected == 0 {
		return types.NewErr(types.NoRowAffectedErr, "no rows affected by the update balance stmt", nil)
	}
	return nil
}

func (r *account) Exists(ctx context.Context, id int64) (bool, error) {
	var exists bool
	err := (*r.txr).GetConn(ctx).QueryRowContext(ctx, "SELECT EXISTS(SELECT id FROM account WHERE id=?)", id).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return exists, types.NewErr(types.SelectStmtErr, "verifying account existence", err)
	}
	return exists, nil
}
