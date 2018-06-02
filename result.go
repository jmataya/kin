package kin

import "database/sql"

// Result contains the results of a database query.
type Result struct {
	Columns []string
	Rows    []*RowResult
}

func newResult(rows *sql.Rows) (*Result, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	res := Result{
		Columns: columns,
		Rows:    []*RowResult{},
	}

	defer rows.Close()
	for rows.Next() {
		row, err := newRowResult(rows)
		if err != nil {
			return nil, err
		}

		res.Rows = append(res.Rows, row)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return &res, nil
}
