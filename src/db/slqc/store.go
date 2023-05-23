package db

import (
	"context"
	"database/sql"
	"desly/util"
	"fmt"
)

type Store interface {
	Querier
	CreateDeslyTx(ctx context.Context, arg CreateDeslyTxParams) (Desly, error)
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

//Store provides all functions to execute db queries and transactions
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		Queries: New(db),
		db: db,
	}
}

//Execute a db transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		//if error rollback db tx
		if rbError := tx.Rollback(); rbError != nil {
			return fmt.Errorf("tx Error %v, rb Error: %v", err, rbError)
		}

		return err
	}

	return tx.Commit()
}

type CreateDeslyTxParams struct {
	Redirect string `json:"redirect"`
}

func (store *SQLStore) CreateDeslyTx(ctx context.Context, arg CreateDeslyTxParams) (Desly, error) {
	var result Desly

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		//get unique desly
		var desly = util.RandomString(6)

		result, err = q.CreateDesly(ctx, CreateDeslyParams{
			Redirect: arg.Redirect,
			Desly:  desly,
		})

		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}