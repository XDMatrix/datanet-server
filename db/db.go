package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type Datastore interface {
	// Ping the database pool
	Ping(ctx context.Context) error
	// Starts a pgx.Tx input context
	BeginTx(ctx context.Context) (pgx.Tx, error)
	// Rollbacks the `pgx.Tx`
	RollbackTx(ctx context.Context, tx pgx.Tx, err error) error
	// Commits the Tx
	CommitTx(ctx context.Context, tx pgx.Tx) error
}

type DBTX interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
}

type PingServicer interface {
	Ping(ctx context.Context, lgr zerolog.Logger) PingResponse
}

type PingResponse struct {
	DBUp bool `json:"db_up"`
}

