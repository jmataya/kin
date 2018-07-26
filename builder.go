package kin

// Builder is a data structure that can be constructed from a database
// row result.
type Builder interface {
	// Columns list the fields that can be extracted from the row result.
	Columns() []FieldBuilder
}

func buildOne(builder Builder, res *RowResult) error {
	for _, column := range builder.Columns() {
		column.Set(res)
	}

	return res.Err()
}
