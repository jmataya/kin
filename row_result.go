package kin

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// RowResult contains the results of a single row in a database result set.
type RowResult struct {
	Columns []string
	Data    map[string]interface{}

	err error
}

// ExtractJSON gets a JSON value from the dataset and unmarshals it into an
// interface passed by the caller.
func (rr *RowResult) ExtractJSON(column string, out interface{}) {
	if rr.err != nil {
		return
	}

	i, ok := rr.Data[column]
	if !ok {
		rr.err = fmt.Errorf("column %s not found in result set", column)
		return
	}

	b, ok := i.(*[]byte)
	if !ok {
		rr.err = fmt.Errorf("column %s (%v) could not be extracted", column, i)
		return
	}

	if err := json.Unmarshal(*b, out); err != nil {
		rr.err = fmt.Errorf("column %s could not be unmarshalled with err %v", column, err)
		return
	}
}

// ExtractInt gets the value in the dataset and returns an integer.
// If the value can't be extracted, it stores an error on the result and
// prevents further extraction from occurring.
func (rr *RowResult) ExtractInt(column string) int {
	if rr.err != nil {
		return 0
	}

	i, ok := rr.Data[column]
	if !ok {
		rr.err = fmt.Errorf("column %s not found in result set", column)
		return 0
	}

	b, ok := i.(*[]byte)
	if !ok {
		rr.err = fmt.Errorf("column %s (%v) could not be extracted as an int", column, i)
		return 0
	}

	num, err := strconv.ParseInt(string(*b), 10, 64)
	if err != nil {
		rr.err = fmt.Errorf("column %s (%+v) could not be extracted as an int with error %v", column, b, err)
		return 0
	}

	return int(num)
}

// ExtractString gets the value in the dataset and returns a string.
// If the value can't be extracted, it stores an error on the result and
// prevents further extraction from occurring.
func (rr *RowResult) ExtractString(column string) string {
	if rr.err != nil {
		return ""
	}

	s, ok := rr.Data[column]
	if !ok {
		rr.err = fmt.Errorf("column %s not found in result set", column)
		return ""
	}

	b, ok := s.(*[]byte)
	if !ok {
		rr.err = fmt.Errorf("column %s (%v) could not be extracted as a string", column, s)
		return ""
	}

	return string(*b)
}

// Err returns any aggregrated errors.
func (rr *RowResult) Err() error {
	return rr.err
}

func newRowResult(row *sql.Rows) (*RowResult, error) {
	columns, err := row.Columns()
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{}

	scanDest := make([]interface{}, len(columns))
	for i := 0; i < len(columns); i++ {
		columnValue := &[]byte{}
		columnName := columns[i]
		data[columnName] = columnValue
		scanDest[i] = columnValue
	}

	if err := row.Scan(scanDest...); err != nil {
		return nil, err
	}

	return &RowResult{
		Columns: columns,
		Data:    data,
		err:     nil,
	}, nil
}

// ExtractTime gets the value in the dataset and returns a time.Time.
// If the value can't be extracted, it stores an error on the result and
// prevents further extraction from occurring.
func (rr *RowResult) ExtractTime(column string) time.Time {
	if rr.err != nil {
		return time.Now()
	}

	i, ok := rr.Data[column]
	if !ok {
		rr.err = fmt.Errorf("column %s not found in result set", column)
		return time.Now()
	}

	b, ok := i.(*[]byte)
	if !ok {
		rr.err = fmt.Errorf("column %s (%v) could not be extracted as an int", column, i)
		return time.Now()
	}

	dateFormat := "2006-01-02T15:04:05.999999Z"
	dateString := string(*b)
	t, err := time.Parse(dateFormat, dateString)
	if err != nil {
		rr.err = fmt.Errorf("column %s (%s) could not be extracted as a timestamp with error %v", column, dateString, err)
		return time.Now()
	}

	return t
}
