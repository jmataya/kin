package kin

import (
	"database/sql"
	"errors"
)

// Transaction is an interface for interacting with a database transaction. It
// abstracts away managing connection pools and various low-level bits.
type Transaction interface {
	// Commit commits the transaction with the database.
	Commit() error

	// Migrate runs migrations.
	Migrate(string) error

	// Rollback the transaction when an error occurs.
	Rollback() error

	// StartTransaction
	StartTransaction() (Transaction, error)

	// Query generates a new query to be executed at a later time.
	Query(stmt string, params ...interface{}) *Query
}

type transaction struct {
	tx *sql.Tx
}

func (t *transaction) Commit() error {
	return t.tx.Commit()
}

func (t *transaction) Migrate(str string) error {
	return errors.New("not implemented")
}

func (t *transaction) Rollback() error {
	return t.tx.Rollback()
}

func (t *transaction) StartTransaction() (Transaction, error) {
	return nil, errors.New("transaction already started")
}

func (t *transaction) Query(stmt string, params ...interface{}) *Query {
	return &Query{
		db:     t.tx,
		stmt:   stmt,
		params: params,
	}
}
