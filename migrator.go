package kin

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

// Migrator is a structure used for running database migrations.
type Migrator struct {
	db *sql.DB
}

// New creates a new Migrator around an existing DB connection.
func NewMigrator(db *sql.DB) (*Migrator, error) {
	if db == nil {
		return nil, errors.New("db connection must be initialized")
	}

	return &Migrator{db: db}, nil
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

	for _, sqlFile := range sqlFiles {
		fmt.Printf("-------- Running %s...", sqlFile)

		file, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", folderPath, sqlFile))
		if err != nil {
			fmt.Printf("FAILED\n")
			fmt.Println("-------- Rolling back changes.")
			return fmt.Errorf("Error executing %s: %v", sqlFile, err)
		}

		_, err = m.db.Exec(string(file))
		if err != nil {
			fmt.Printf("FAILED\n")
			fmt.Println("-------- Rolling back changes.")
			return fmt.Errorf("Error executing %s: %v", sqlFile, err)
		}

		fmt.Printf("COMPLETED\n")
	}

	return nil
}

func fileSuffix(fileName string) string {
	parts := strings.Split(fileName, ".")
	return parts[len(parts)-1]
}
