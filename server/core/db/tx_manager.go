package db

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TxManager[T any] interface {
	ExecTx(context.Context, func(queries T) error) error
}

type PgTxManager[T any] struct {
	pool       *pgxpool.Pool
	newQuerier func(pgx.Tx) T
}

func (tm *PgTxManager[T]) ExecTx(ctx context.Context, fn func(queries T) error) error {
	tx, err := tm.pool.Begin(ctx)
	if err != nil {
		return err
	}
	var finalErr error
	defer func() {
		if finalErr != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				finalErr = rbErr
			}
			return
		}
		if cmErr := tx.Commit(ctx); cmErr != nil {
			finalErr = cmErr
		}
	}()

	querier := tm.newQuerier(tx)
	finalErr = fn(querier)

	return finalErr
}

func NewPgTxManager[T any](pool *pgxpool.Pool, newQuerier func(pgx.Tx) T) TxManager[T] {
	return &PgTxManager[T]{
		pool:       pool,
		newQuerier: newQuerier,
	}
}
