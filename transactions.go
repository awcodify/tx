package tx

import (
	"context"
	"database/sql"
)

type Tx struct {
	db *sql.DB
}

type Transaction interface {
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Prepare(query string) (*sql.Stmt, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

// Initialize Tx with DB connection
func NewTx(db *sql.DB) Tx {
	return Tx{db: db}
}

// Wrap the transaction.
func (tx Tx) Wrap(fn func(t Transaction) error) (err error) {
	// Starts the transaction
	t, err := db.Begin()
	if err != nil {
		return
	}

	// Push the Commit() and Rollback() onto the list
	//	Everythings well, then commit to the database.
	//	If an error returned, it should be Rollback()
	defer func() {
		switch err {
		case nil:
			err = t.Commit()
		default:
			t.Rollback()
		}
	}()

	// Run the wrapped function
	err = fn(t)

	return
}
