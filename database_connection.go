package kin

import "database/sql"

type databaseConnection interface {
	Prepare(string) (*sql.Stmt, error)
}
