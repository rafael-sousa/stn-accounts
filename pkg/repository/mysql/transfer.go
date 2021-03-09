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
	transfers := make([]*entity.Transfer, 0)
	for rows.Next() {
		transfer := entity.Transfer{}
		err = rows.Scan(&transfer.ID, &transfer.Origin, &transfer.Destination, &transfer.Amount, &transfer.CreatedAt)
		if err != nil {
			return nil, types.NewErr(types.SelectStmtErr, "scanning the transfer row", &err)
		}
		transfers = append(transfers, &transfer)
	}
	if err = rows.Err(); err != nil {
		return nil, types.NewErr(types.SelectStmtErr, "iterating over the transfer rows", &err)
	}
	return transfers, nil

}
func (r *transfer) Create(ctx context.Context, transfer *entity.Transfer) (*entity.Transfer, error) {
	stmt, err := (*r.txr).GetConn(ctx).PrepareContext(ctx, "INSERT INTO transfer(account_origin_id, account_destination_id, amount, created_at) VALUES (?,?,?,?)")
	if err != nil {
		return nil, types.NewErr(types.SelectStmtErr, "preparing transfer insert stmt", &err)
	}
	defer stmt.Close()
	result, err := stmt.ExecContext(ctx, transfer.Origin, transfer.Destination, transfer.Amount, transfer.CreatedAt)
	if err != nil {
		return nil, types.NewErr(types.SelectStmtErr, "exec transfer insert stmt", &err)
	}
	if transfer.ID, err = result.LastInsertId(); err != nil {
		return nil, types.NewErr(types.SelectStmtErr, "getting the inserted transfer id", &err)
	}
	return transfer, nil
}

// NewTransfer creates a value that satisfies the repository.Transfer interface
func NewTransfer(txr *repository.Transactioner) repository.Transfer {
	return &transfer{txr: txr}
}
