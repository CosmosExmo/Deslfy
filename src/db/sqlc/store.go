package db

import (
	"context"
	"database/sql"
)

type Store interface {
	Querier
	CreateDeslyTx(ctx context.Context, arg CreateDeslyTxParams) (Desly, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

// Store provides all functions to execute db queries and transactions
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}

