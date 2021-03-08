package kin

// Model is a data structure that can be used to construct a database
// insert or update.
type Model interface {
	// Columns list the fields that can be extracted from the row result.
	Columns() []FieldBuilder

	// TableName returns the name of the table that this model maps to.
	TableName() string
}
