package kin

import "time"

// FieldBuilder is an interface for building up fields from a database
// row result.
type FieldBuilder interface {
	// Set assigns the wrapped field with the result from a row.
	Set(res *RowResult)
}

// BoolField is a reference to an integer field.
type BoolField struct {
	fieldName string
	field     *bool
}

// Set assigns the wrapped field with the result from a row.
func (b BoolField) Set(res *RowResult) {
	*b.field = res.ExtractBool(b.fieldName)
}

// IntField is a reference to an integer field.
type IntField struct {
	fieldName string
	field     *int
}

// Set assigns the wrapped field with the result from a row.
func (i IntField) Set(res *RowResult) {
	*i.field = res.ExtractInt(i.fieldName)
}

// JSONField is a reference to a string field.
type JSONField struct {
	fieldName string
	field     interface{}
}

// Set assigns the wrapped field with the result from a row.
func (j JSONField) Set(res *RowResult) {
	res.ExtractJSON(j.fieldName, j.field)
}

// StringField is a reference to a string field.
type StringField struct {
	fieldName string
	field     *string
}

// Set assigns the wrapped field with the result from a row.
func (s StringField) Set(res *RowResult) {
	*s.field = res.ExtractString(s.fieldName)
}

// TimeField is a reference to a string field.
type TimeField struct {
	fieldName string
	field     *time.Time
}

// Set assigns the wrapped field with the result from a row.
func (t TimeField) Set(res *RowResult) {
	*t.field = res.ExtractTime(t.fieldName)
}
