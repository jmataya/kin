package kin

import (
	"database/sql"
	"errors"
)

// Database is an interface for interacting with the database. It abstracts
// away managing connection pools and various other low-level bits.
type Database interface {
	// Query generates a new query to be executed at a later time.
	Query(stmt string, params ...interface{}) *Query

	// StartTransaction initiates a database transaction object.
	StartTransaction() (Transaction, error)
}

// New creates a new wrapper around an existing DB connection.
func New(db *sql.DB) (Database, error) {
	if db == nil {
		return nil, errors.New("db connection must be initialized")
	}

	return &database{db: db}, nil
}

type database struct {
	db *sql.DB
}

func (d *database) Query(stmt string, params ...interface{}) *Query {
	return &Query{
		db:     d.db,
		stmt:   stmt,
		params: params,
	}
}

func (d *database) StartTransaction() (Transaction, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
	}

	return &transaction{tx: tx}, nil
}
