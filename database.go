package kin

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq" // Needed to initialize the Postgres SQL driver.
)

// Database is an interface for interacting with the database. It abstracts
// away managing connection pools and various other low-level bits.
type Database interface {
	// Close terminates the database connection.
	Close() error

	// Insert generates an insert query for a model.
	Insert(m Model) *Query

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

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("unable to connect to database %v", err)
	}

	return &database{db: db}, nil
}

// NewConnection initializes a new connection and creates a wrapper around it.
func NewConnection(dbURL string) (Database, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection %v", err)
	}

	return New(db)
}

type database struct {
	db *sql.DB
}

func (d *database) Close() error {
	return d.db.Close()
}

func (d *database) Insert(m Model) *Query {
	var columns string
	var values string
	var params []interface{}

	for _, column := range m.Columns() {
		if !column.IsSet() {
			continue
		}

		separator := ""
		if len(params) > 0 {
			separator = ", "
		}

		params = append(params, column.Get())
		columns = fmt.Sprintf("%s%s%s", columns, separator, column.FieldName())
		values = fmt.Sprintf("%s%s$%d", values, separator, len(params))
	}

	stmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING *", m.TableName(), columns, values)

	return d.Query(stmt, params...)
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
