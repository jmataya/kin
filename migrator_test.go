package kin

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func setupMigrationDir(path string) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func cleanupMigrationDir(path string) error {
	return os.RemoveAll(path)
}

func createFile(path, name, contents string) error {
	file, err := os.Create(fmt.Sprintf("%s/%s", path, name))
	if err != nil {
		return err
	}

	defer file.Close()

	if _, err := file.WriteString(fmt.Sprintf("%s\n", contents)); err != nil {
		return err
	}

	file.Sync()
	return nil
}

func TestMigrate(t *testing.T) {
	migrationPath := "./sql"

	cleanupMigrationDir(migrationPath)
	if err := setupMigrationDir(migrationPath); err != nil {
		t.Errorf("setupMigrationDir = %v", err)
		return
	}

	createFoo := "create table foo (id serial primary key);"
	createBar := `
  	create table bar (
  		id serial primary key,
  		foo_id int not null,
  		foreign key (foo_id) references foo(id) on update restrict on delete restrict
  	);
  `
	if err := createFile(migrationPath, "1__create_foo.sql", createFoo); err != nil {
		t.Errorf("createFile = %v", err)
		cleanupMigrationDir(migrationPath)
		return
	}

	if err := createFile(migrationPath, "2__create_bar.sql", createBar); err != nil {
		t.Errorf("createFile = %v", err)
		cleanupMigrationDir(migrationPath)
		return
	}

	connStr := os.Getenv("POSTGRES_URL")
	if connStr == "" {
		panic("POSTGRES_URL must not be empty")
	}

	db, _ := sql.Open("postgres", connStr)
	migrator, _ := NewMigrator(db)

	if err := migrator.Migrate(migrationPath); err != nil {
		t.Errorf("migrator.Migrate(%s) = %v, want nil", migrationPath, err)
		cleanupMigrationDir(migrationPath)
	}

	cleanupMigrationDir(migrationPath)
}

func TestMigrateRerun(t *testing.T) {
	migrationPath := "./sql"

	cleanupMigrationDir(migrationPath)
	if err := setupMigrationDir(migrationPath); err != nil {
		t.Errorf("setupMigrationDir = %v", err)
		return
	}

	createBaz := "create table baz (id serial primary key);"
	if err := createFile(migrationPath, "1__create_baz.sql", createBaz); err != nil {
		t.Errorf("createFile = %v", err)
		cleanupMigrationDir(migrationPath)
		return
	}

	connStr := os.Getenv("POSTGRES_URL")
	if connStr == "" {
		panic("POSTGRES_URL must not be empty")
	}
	db, _ := sql.Open("postgres", connStr)
	migrator, _ := NewMigrator(db)

	if err := migrator.Migrate(migrationPath); err != nil {
		t.Errorf("migrator.Migrate(%s) = %v, want nil", migrationPath, err)
		cleanupMigrationDir(migrationPath)
		return
	}

	if err := migrator.Migrate(migrationPath); err != nil {
		t.Errorf("migrator.Migrate(%s) = %v, want nil", migrationPath, err)
	}

	cleanupMigrationDir(migrationPath)
}
