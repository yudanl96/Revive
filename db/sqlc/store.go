package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
}

// extend to include transactions in sql database
type SQLStore struct {
	*Queries
	db *sql.DB //create db transaction
}

// new instance
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execute a function
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil) //using default isolation level
	if err != nil {
		return err
	}

	queries := New(tx)
	err = fn(queries)
	if err != nil { // if error rollback
		if rBErr := tx.Rollback(); rBErr != nil {
			return fmt.Errorf("transaction err: %v, rollback err: %v", err, rBErr)
		}
		return err
	}

	return tx.Commit()
}

// some transactions to implement
