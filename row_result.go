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
	rawCol, err := rr.extractColumn(column)
	if err != nil {
		rr.err = err
		return
	}

	if err := json.Unmarshal(rawCol, out); err != nil {
		rr.err = fmt.Errorf("column %s could not be unmarshalled with err %v", column, err)
		return
	}
}

// ExtractBool gets the value in the dataset and returns a boolean.
// If the value can't be extracted, it stores an error on the result and
// prevents further extraction from occurring.
func (rr *RowResult) ExtractBool(column string) bool {
	rawCol, err := rr.extractColumn(column)
	if err != nil {
		rr.err = err
		return false
	}

	boolVal, err := strconv.ParseBool(string(rawCol))
	if err != nil {
		rr.err = fmt.Errorf("column %s (%+v) could not be extracted as a bool with error %v", column, rawCol, err)
		return false
	}

	return boolVal
}

// ExtractInt gets the value in the dataset and returns an integer.
// If the value can't be extracted, it stores an error on the result and
// prevents further extraction from occurring.
func (rr *RowResult) ExtractInt(column string) int {
	rawCol, err := rr.extractColumn(column)
	if err != nil {
		rr.err = err
		return 0
	}

	num, err := strconv.ParseInt(string(rawCol), 10, 64)
	if err != nil {
		rr.err = fmt.Errorf("column %s (%+v) could not be extracted as an int with error %v", column, rawCol, err)
		return 0
	}

	return int(num)
}

// ExtractString gets the value in the dataset and returns a string.
// If the value can't be extracted, it stores an error on the result and
// prevents further extraction from occurring.
func (rr *RowResult) ExtractString(column string) string {
	rawCol, err := rr.extractColumn(column)
	if err != nil {
		rr.err = err
		return ""
	}

	return string(rawCol)
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
	rawCol, err := rr.extractColumn(column)
	if err != nil {
		rr.err = err
		return time.Now()
	}

	dateFormat := "2006-01-02T15:04:05.999999Z"
	dateString := string(rawCol)
	t, err := time.Parse(dateFormat, dateString)
	if err != nil {
		rr.err = fmt.Errorf("column %s (%s) could not be extracted as a timestamp with error %v", column, dateString, err)
		return time.Now()
	}

	return t
}

func (rr *RowResult) extractColumn(column string) ([]byte, error) {
	if rr.err != nil {
		return nil, rr.err
	}

	col, ok := rr.Data[column]
	if !ok {
		return nil, fmt.Errorf("column %s not found in result set", column)
	}

	b, ok := col.(*[]byte)
	if !ok {
		return nil, fmt.Errorf("column %s (%v) could not be extracted", column, col)
	}

	return *b, nil
}
