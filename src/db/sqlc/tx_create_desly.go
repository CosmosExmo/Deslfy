package db

import "context"

type CreateDeslyTxParams struct {
	Redirect string `json:"redirect"`
	Owner    string `json:"owner"`
}

func (store *SQLStore) CreateDeslyTx(ctx context.Context, arg CreateDeslyTxParams) (Desly, error) {
	var result Desly

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		params := CreateDeslyParams(arg)

		result, err = q.CreateDesly(ctx, params)

		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}