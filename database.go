package kin

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

// Database is an interface for interacting with the database. It abstracts
// away managing connection pools and various other low-level bits.
type Database interface {
	// Migrate runs a series of database migrations.
	Migrate(folderPath string) error

	// Query generates a new query to be executed at a later time.
	Query(stmt string, params ...interface{}) *Query

	// StartTransaction initiates a database transaction object.
	StartTransaction() (Transaction, error)
}

// New creates a new wrapper around an existing DB connection.
func New(db *sql.DB) (Database, error) {
	if db == nil {
		return nil, errors.New("db connection must be initialized")
	}

	return &database{db: db}, nil
}

type database struct {
	db *sql.DB
}

func fileSuffix(fileName string) string {
	parts := strings.Split(fileName, ".")
	return parts[len(parts)-1]
}

func (d *database) Migrate(folderPath string) error {
	fmt.Println("Starting database migrations...")
	fmt.Println("")

	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return fmt.Errorf("unable to read migrations: %v", err)
	}

	sqlFiles := []string{}
	for _, file := range files {
		if !file.IsDir() && fileSuffix(file.Name()) == "sql" {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}

	if len(sqlFiles) == 0 {
		return fmt.Errorf("no migrations found in folder '%s'", folderPath)
	}

	for _, sqlFile := range sqlFiles {
		fmt.Printf("-------- Running %s...", sqlFile)

		file, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", folderPath, sqlFile))
		if err != nil {
			fmt.Printf("FAILED\n")
			fmt.Println("-------- Rolling back changes.")
			return fmt.Errorf("Error executing %s: %v", sqlFile, err)
		}

		_, err = d.db.Exec(string(file))
		if err != nil {
			fmt.Printf("FAILED\n")
			fmt.Println("-------- Rolling back changes.")
			return fmt.Errorf("Error executing %s: %v", sqlFile, err)
		}

		fmt.Printf("COMPLETED\n")
	}
	return nil
}

func (d *database) Query(stmt string, params ...interface{}) *Query {
	return &Query{
		db:     d.db,
		stmt:   stmt,
		params: params,
	}
}

func (d *database) StartTransaction() (Transaction, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
	}

	return &transaction{tx: tx}, nil
}
