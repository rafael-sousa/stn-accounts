package repository

import (
	"context"
	"database/sql"

	"github.com/rafael-sousa/stn-accounts/pkg/model/types"
	"github.com/rs/zerolog/log"
)

type txKey string

// CtxTxKey is the context key that holds a tx value in a context.Context
const CtxTxKey txKey = "txKey"

type transactioner struct {
	db *sql.DB
}

// Transactioner handles operations that store transactions in contexts and run functions within a transactional concept
type Transactioner interface {
	WithTx(ctx context.Context, f func(context.Context) error) (err error)
	GetConn(ctx context.Context) Connection
}

// WithTx starts a db transaction and stores it on the specified context.
// It runs the function stored at f with the transactional context.
// If f function yields no error, the transaction is committed otherwise is rolled back
func (manager *transactioner) WithTx(ctx context.Context, f func(context.Context) error) (err error) {
	ctxTx := ctx.Value(CtxTxKey)

	if ctxTx == nil {
		ctxTx, err := manager.db.BeginTx(ctx, nil)
		if err != nil {
			return types.NewErr(types.InternalErr, "unable to begin tx", &err)
		}
		ctx = context.WithValue(ctx, CtxTxKey, ctxTx)
	}

	err = f(ctx)

	if tx, ok := ctxTx.(*sql.Tx); ok {
		if err == nil {
			return tx.Commit()
		}
		if e := tx.Rollback(); e != nil {
			log.Error().Caller().Err(e).Msg("unable to rollback tx")
		}
	}
	return err
}

// GetConn returns a connection to the current database pool
func (manager *transactioner) GetConn(ctx context.Context) Connection {
	if k, ok := ctx.Value(CtxTxKey).(Connection); ok {
		return k
	}
	return manager.db
}

// NewTxr creates a new Transactioner value
func NewTxr(db *sql.DB) Transactioner {
	return &transactioner{db: db}
}
