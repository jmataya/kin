package kin

import "time"

// FieldBuilder is an interface for building up fields from a database
// row result.
type FieldBuilder interface {
	// Set assigns the wrapped field with the result from a row.
	Set(res *RowResult)
}

// BoolField creates a reference to a boolean field.
func BoolField(fieldName string, field *bool) FieldBuilder {
	return boolField{fieldName, field}
}

type boolField struct {
	fieldName string
	field     *bool
}

func (b boolField) Set(res *RowResult) {
	*b.field = res.ExtractBool(b.fieldName)
}

// DecimalField creates a reference to a field for a floating point number.
func DecimalField(fieldName string, field *float64) FieldBuilder {
	return decimalField{fieldName, field}
}

type decimalField struct {
	fieldName string
	field     *float64
}

func (d decimalField) Set(res *RowResult) {
	*d.field = res.ExtractDecimal(d.fieldName)
}

// IntField creates a reference to an integer field.
func IntField(fieldName string, field *int) FieldBuilder {
	return intField{fieldName, field}
}

type intField struct {
	fieldName string
	field     *int
}

func (i intField) Set(res *RowResult) {
	*i.field = res.ExtractInt(i.fieldName)
}

// JSONField creates a reference to a JSON field.
func JSONField(fieldName string, field interface{}) FieldBuilder {
	return jsonField{fieldName, field}
}

type jsonField struct {
	fieldName string
	field     interface{}
}

func (j jsonField) Set(res *RowResult) {
	res.ExtractJSON(j.fieldName, j.field)
}

// StringField creates a reference to a string field.
func StringField(fieldName string, field *string) FieldBuilder {
	return stringField{fieldName, field}
}

type stringField struct {
	fieldName string
	field     *string
}

func (s stringField) Set(res *RowResult) {
	*s.field = res.ExtractString(s.fieldName)
}

// TimeField creates a reference to a string field.
func TimeField(fieldName string, field *time.Time) FieldBuilder {
	return timeField{fieldName, field}
}

type timeField struct {
	fieldName string
	field     *time.Time
}

func (t timeField) Set(res *RowResult) {
	*t.field = res.ExtractTime(t.fieldName)
}
