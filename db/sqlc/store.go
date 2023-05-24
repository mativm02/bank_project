package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute database queries.
type Store interface {
	TransferTx(ctx context.Context, arg CreateTransferParams) (TransferTxResults, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	Querier
}

// Store provides all functions to execute database SQL queries.
type SQLStore struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new store.
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
