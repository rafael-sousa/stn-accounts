package mysql

import (
	"context"

	"github.com/rafael-sousa/stn-accounts/pkg/model/entity"
	"github.com/rafael-sousa/stn-accounts/pkg/model/types"
	"github.com/rafael-sousa/stn-accounts/pkg/repository"
)

type transfer struct {
	txr *repository.Transactioner
}

var _ repository.Transfer = (*transfer)(nil)

func (r *transfer) Fetch(ctx context.Context, origin int64) ([]*entity.Transfer, error) {
	rows, err := (*r.txr).GetConn(ctx).QueryContext(ctx, "SELECT id, account_origin_id, account_destination_id, amount, created_at FROM transfer WHERE account_origin_id=?", origin)
	if err != nil {
		return nil, types.NewErr(types.SelectStmtErr, "querying transfers by id", &err)
	}
	defer rows.Close()
	tranfers := make([]*entity.Transfer, 0)
	for rows.Next() {
		e := entity.Transfer{}
		err = rows.Scan(&e.ID, &e.Origin, &e.Destination, &e.Amount, &e.CreatedAt)
		if err != nil {
			return nil, types.NewErr(types.SelectStmtErr, "scanning the transfer row", &err)
		}
		tranfers = append(tranfers, &e)
	}
	err = rows.Err()
	if err != nil {
		return nil, types.NewErr(types.SelectStmtErr, "iterating over the transfer rows", &err)
	}
	return tranfers, nil

}
func (r *transfer) Create(ctx context.Context, e *entity.Transfer) (*entity.Transfer, error) {
	stmt, err := (*r.txr).GetConn(ctx).PrepareContext(ctx, "INSERT INTO transfer(account_origin_id, account_destination_id, amount, created_at) VALUES (?,?,?,?)")

	if err != nil {
		return nil, types.NewErr(types.SelectStmtErr, "preparing transfer insert stmt", &err)
	}
	defer stmt.Close()
	result, err := stmt.ExecContext(ctx, e.Origin, e.Destination, e.Amount, e.CreatedAt)
	if err != nil {
		return nil, types.NewErr(types.SelectStmtErr, "exec transfer insert stmt", &err)
	}
	e.ID, err = result.LastInsertId()
	if err != nil {
		return nil, types.NewErr(types.SelectStmtErr, "getting the inserted transfer id", &err)
	}
	return e, nil
}

// NewTransfer creates a value that satisfies the repository.Transfer interface
func NewTransfer(txr *repository.Transactioner) repository.Transfer {
	return &transfer{txr: txr}
}
