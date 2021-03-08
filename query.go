package kin

import (
	"database/sql"
	"errors"
)

type databaseConnection interface {
	Prepare(string) (*sql.Stmt, error)
}

// Query is a SQL query that has yet to be executed.
type Query struct {
	db     databaseConnection
	stmt   string
	params []interface{}
}

// One executes the query and returns an error if no results are found.
func (q Query) One() (*RowResult, error) {
	result, err := q.Run()
	if err != nil {
		return nil, err
	}

	if len(result.Rows) < 1 {
		return nil, errors.New("query must return at least one result")
	}

	return result.Rows[0], nil
}

// OneAndExtract executes the query and updates the builder with the result.
func (q Query) OneAndExtract(b Builder) error {
	res, err := q.One()
	if err != nil {
		return err
	}

	return buildOne(b, res)
}

// OneAndExtractFn executes the query and returns the model with the result
// with an extraction function.
func (q Query) OneAndExtractFn(buildFn func(*RowResult) error) error {
	res, err := q.One()
	if err != nil {
		return err
	}

	return buildFn(res)
}

// ExtractFn executes the query and iterates over all rows to extract
// the result.
func (q Query) ExtractFn(buildFn func(*RowResult) error) error {
	result, err := q.Run()
	if err != nil {
		return err
	}

	for _, row := range result.Rows {
		if err := buildFn(row); err != nil {
			return err
		}
	}

	return nil
}

// Run executes the query and returns the results.
func (q Query) Run() (*Result, error) {
	stmt, err := q.db.Prepare(q.stmt)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(q.params...)
	if err != nil {
		return nil, err
	}

	return newResult(rows)
}
