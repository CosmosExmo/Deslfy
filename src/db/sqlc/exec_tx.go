package db

import (
	"context"
	"fmt"
)

// ExecTx executes a function within a database transaction
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