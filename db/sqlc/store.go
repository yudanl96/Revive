package db

import (
	"context"
	"database/sql"
	"fmt"
)

// extend to include transactions, composition is preferred than inheritance
type Store struct {
	*Queries
	db *sql.DB //create db transaction
}

// new instance
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execute a function
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
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
