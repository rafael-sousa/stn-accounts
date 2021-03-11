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

// Transactioner handles operations that store transactions in contexts and run functions within a transactional concept
type Transactioner interface {
	WithTx(ctx context.Context, fn func(context.Context) error) (err error)
	GetConn(ctx context.Context) Connection
}

type transactioner struct {
	db *sql.DB
}

// NewTxr creates a new Transactioner value
func NewTxr(db *sql.DB) Transactioner {
	return &transactioner{db: db}
}

// WithTx starts a db transaction and stores it on the specified context.
// It runs the function stored at fn with the transactional context.
// If f function yields no error, the transaction is committed otherwise is rolled back
func (txr *transactioner) WithTx(ctx context.Context, fn func(context.Context) error) (err error) {
	ctxTx := ctx.Value(CtxTxKey)

	if ctxTx == nil {
		ctxTx, err = txr.db.BeginTx(ctx, nil)
		if err != nil {
			return types.NewErr(types.InternalErr, "unable to begin tx", err)
		}
		ctx = context.WithValue(ctx, CtxTxKey, ctxTx)
	}

	err = fn(ctx)

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
func (txr *transactioner) GetConn(ctx context.Context) Connection {
	if conn, ok := ctx.Value(CtxTxKey).(Connection); ok {
		return conn
	}
	return txr.db
}
