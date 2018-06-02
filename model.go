package kin

// Model is a standard way to interface with the database.
type Model interface {
	Build(*RowResult) error
}
