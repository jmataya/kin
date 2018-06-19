package kin

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"strings"
)

const (
	sqlCreateSchemasTable = `
		CREATE TABLE IF NOT EXISTS schemas (
			id serial primary key,
			filename text not null check(length(filename) <= 255),
			applied_on timestamp without time zone default (now() at time zone 'utc')
		)
	`

	sqlInsertSchema = "INSERT INTO schemas (filename) VALUES ($1)"
)

// Migrator is a structure used for running database migrations.
type Migrator struct {
	db Database
}

// New creates a new Migrator around an existing DB connection.
func NewMigrator(db *sql.DB) (*Migrator, error) {
	dbi, err := New(db)
	if err != nil {
		return nil, err
	}

	return &Migrator{db: dbi}, nil
}

// Migrate runs a series of database migrations.
func (m Migrator) Migrate(folderPath string) error {
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

	txn, err := m.db.StartTransaction()
	if err != nil {
		return fmt.Errorf("Unexpected error starting transaction: %v", err)
	}

	fmt.Printf("-------- Ensuring database is set up...")
	if err := txn.Exec(sqlCreateSchemasTable); err != nil {
		fmt.Printf("FAILED\n")
		return fmt.Errorf("Error setting up schemas table: %v", err)
	}

	res, err := txn.Query("SELECT * FROM schemas").Run()
	if err != nil {
		txn.Rollback()
		return fmt.Errorf("Unable to get applied migrations: %v", err)
	}

	appliedMigrations := map[string]int{}
	for _, rowRes := range res.Rows {
		id := rowRes.ExtractInt("id")
		filename := rowRes.ExtractString("filename")
		appliedMigrations[filename] = id
	}

	for _, sqlFile := range sqlFiles {
		fmt.Printf("-------- Running %s...", sqlFile)
		if _, ok := appliedMigrations[sqlFile]; ok {
			fmt.Printf("SKIPPED\n")
			continue
		}

		file, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", folderPath, sqlFile))
		if err != nil {
			fmt.Printf("FAILED\n")
			fmt.Println("-------- Rolling back changes.")
			txn.Rollback()
			return fmt.Errorf("Error executing %s: %v", sqlFile, err)
		}

		err = txn.Exec(string(file))
		if err != nil {
			fmt.Printf("FAILED\n")
			fmt.Println("-------- Rolling back changes.")
			txn.Rollback()
			return fmt.Errorf("Error executing %s: %v", sqlFile, err)
		}

		err = txn.Exec(sqlInsertSchema, sqlFile)
		if err != nil {
			fmt.Printf("FAILED\n")
			fmt.Println("-------- Rolling back changes.")
			txn.Rollback()
			return fmt.Errorf("Error updating schemas table with %s: %v", sqlFile, err)
		}

		fmt.Printf("COMPLETED\n")
	}

	return txn.Commit()
}

func fileSuffix(fileName string) string {
	parts := strings.Split(fileName, ".")
	return parts[len(parts)-1]
}
