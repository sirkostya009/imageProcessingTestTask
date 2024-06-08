package db

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (q *Queries) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return q.db.(*pgxpool.Pool).Begin(ctx)
}
