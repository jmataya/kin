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

	// Exec runs a query against the database that doesn't return any results.
	Exec(query string, args ...interface{}) error

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

func (t *transaction) Exec(query string, args ...interface{}) error {
	_, err := t.tx.Exec(query, args...)
	return err
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
