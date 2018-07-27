package kin

import (
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

const (
	createBuilderTable = `
		create table builders (
			id serial primary key,
			name text not null,
			attributes jsonb not null default '{}',
			is_active boolean not null default false,
			created_at timestamp without time zone
		);
	`

	sqlInsertBuilder = `
		INSERT INTO builders (name, attributes, is_active, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING *
	`
)

type builderModel struct {
	ID         int
	Name       string
	Attributes map[string]interface{}
	IsActive   bool
	CreatedAt  time.Time
}

func (b *builderModel) Columns() []FieldBuilder {
	return []FieldBuilder{
		IntField{"id", &b.ID},
		StringField{"name", &b.Name},
		JSONField{"attributes", &b.Attributes},
		BoolField{"is_active", &b.IsActive},
		TimeField{"created_at", &b.CreatedAt},
	}
}

func TestBuild(t *testing.T) {
	migrationPath := "./sql"
	cleanupMigrationDir(migrationPath)
	if err := setupMigrationDir(migrationPath); err != nil {
		t.Errorf("setupMigrationDir = %v", err)
		return
	}

	if err := createFile(migrationPath, "1__create_builders.sql", createBuilderTable); err != nil {
		t.Errorf("createFile = %v", err)
		cleanupMigrationDir(migrationPath)
		return
	}

	connStr := os.Getenv("POSTGRES_URL")
	if connStr == "" {
		panic("POSTGRES_URL must not be empty")
	}

	sqlDB, _ := sql.Open("postgres", connStr)
	defer sqlDB.Close()

	migrator, _ := NewMigrator(sqlDB)

	if err := migrator.Migrate(migrationPath); err != nil {
		t.Errorf("migrator.Migrate(%s) = %v, want nil", migrationPath, err)
		cleanupMigrationDir(migrationPath)
		return
	}

	name := "Builder Test"
	attributes := `{ "lang": "en" }`
	isActive := true
	createdAt := time.Now().UTC()

	builder := new(builderModel)

	db, _ := New(sqlDB)
	err := db.Query(sqlInsertBuilder, name, attributes, isActive, createdAt).OneAndExtract(builder)
	if err != nil {
		t.Errorf("db.Query(...).OneAndExtract(...) = %v, want <nil>", err)
		return
	}

	if builder.ID == 0 {
		t.Error("builder.ID = 0, want > 0")
	}

	if builder.Name != name {
		t.Errorf("builder.Name = %v, want %v", builder.Name, name)
	}

	if builder.IsActive != isActive {
		t.Errorf("builder.IsActive = %v, want %v", builder.IsActive, isActive)
	}

	if builder.CreatedAt.Equal(createdAt) {
		t.Errorf("builder.CreatedAt = %v, want %v", builder.CreatedAt, createdAt)
	}

	if builder.Attributes["lang"] != "en" {
		t.Error(`"builder.Attributes != { "lang": "en" }`)
	}

	cleanupMigrationDir(migrationPath)
}