package kin

import (
	"database/sql"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/jmataya/renv/autoload"
	_ "github.com/lib/pq"
)

func TestDatabaseConnect(t *testing.T) {
	connStr := os.Getenv("POSTGRES_URL")
	if connStr == "" {
		t.Error("POSTGRES_URL is empty")
		return
	}

	sqlDB, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Errorf("sql.Open(...) = %v, want <nil>", err)
		return
	}

	if _, err := New(sqlDB); err != nil {
		t.Errorf("New(...) = (_, %v), want (_, <nil>)", err)
	}
}

func TestDatabaseConnectURL(t *testing.T) {
	connStr := os.Getenv("POSTGRES_URL")
	if connStr == "" {
		t.Error("POSTGRES_URL is empty")
		return
	}

	if _, err := NewConnection(connStr); err != nil {
		t.Errorf("NewConnection(%s) = %v, want <nil>", connStr, err)
	}
}

func TestDatabaseNoConnection(t *testing.T) {
	want := "unable to create connection"
	if _, err := New(nil); err == nil {
		t.Errorf("New(...) = (_, <nil>), want (_, %s)", want)
	} else if err.Error() == want {
		t.Errorf("New(...) = (_, %v), want (_, %s)", err, want)
	}
}

func TestDatabaseBadConnection(t *testing.T) {
	want := "unable to create connection"
	sqlDB, _ := sql.Open("postgres", "")

	if _, err := New(sqlDB); err == nil {
		t.Errorf("New(...) = (_, <nil>), want (_, %s)", want)
	} else if err.Error() == want {
		t.Errorf("New(...) = (_, %v), want (_, %s)", err, want)
	}
}

func TestDatabaseBadConnectionURL(t *testing.T) {
	want := "unable to create connection"
	if _, err := NewConnection(""); err == nil {
		t.Errorf("NewConnection(\"\") = (_, <nil>), want (_, %s)", want)
	} else if strings.HasPrefix(err.Error(), want) {
		t.Errorf("NewConnection(\"\") = (_, %v), want (_, %s)", err, want)
	}
}

func TestDatabaseTransactionError(t *testing.T) {
	connStr := os.Getenv("POSTGRES_URL")
	if connStr == "" {
		t.Error("POSTGRES_URL is empty")
		return
	}

	db, _ := NewConnection(connStr)
	db.Close()

	want := "sql: database is closed"
	if _, err := db.StartTransaction(); err == nil {
		t.Errorf("NewConnection(%s) = <nil>, want %s", connStr, want)
	} else if err.Error() != want {
		t.Errorf("NewConnection(%s) = %v, want %s", connStr, err, want)
	}
}

type testModel struct {
	ID   int
	Name string
}

func (tm testModel) TableName() string {
	return "test_model"
}

func (tm testModel) Columns() []FieldBuilder {
	return []FieldBuilder{
		IntField("id", &tm.ID),
		StringField("name", &tm.Name),
	}
}

func TestInsert(t *testing.T) {
	tm := testModel{ID: 1, Name: "Donkey Hote"}
	db := &database{}
	q := db.Insert(tm)

	wantStmt := "INSERT INTO test_model (id, name) VALUES ($1, $2) RETURNING *"
	if q.stmt != wantStmt {
		t.Errorf("q.stmt = %s, want %s", q.stmt, wantStmt)
	}

	wantParams := []interface{}{1, "Donkey Hote"}
	if len(q.params) != len(wantParams) {
		t.Errorf("len(q.params) = %d, want %d", len(q.params), len(wantParams))
		return
	}

	for i, param := range q.params {
		if param != wantParams[i] {
			t.Errorf("q.params[%d] = %v, want %v", i, param, wantParams[i])
		}
	}
}

type testModelTime struct {
	ID        int
	Name      string
	CreatedAt time.Time
}

func (tm testModelTime) TableName() string {
	return "test_model"
}

func (tm testModelTime) Columns() []FieldBuilder {
	return []FieldBuilder{
		IntField("id", &tm.ID),
		StringField("name", &tm.Name),
		TimeField("created_at", &tm.CreatedAt),
	}
}

func TestInsertUnsetTime(t *testing.T) {
	tm := testModelTime{Name: "Donkey Hote"}
	db := &database{}
	q := db.Insert(tm)

	wantStmt := "INSERT INTO test_model (name) VALUES ($1) RETURNING *"
	if q.stmt != wantStmt {
		t.Errorf("q.stmt = %s, want %s", q.stmt, wantStmt)
	}

	wantParams := []interface{}{"Donkey Hote"}
	if len(q.params) != len(wantParams) {
		t.Errorf("len(q.params) = %d, want %d", len(q.params), len(wantParams))
		return
	}

	for i, param := range q.params {
		if param != wantParams[i] {
			t.Errorf("q.params[%d] = %v, want %v", i, param, wantParams[i])
		}
	}
}

func TestInsertSetTime(t *testing.T) {
	createdAt := time.Now()
	tm := testModelTime{Name: "Donkey Hote", CreatedAt: createdAt}
	db := &database{}
	q := db.Insert(tm)

	wantStmt := "INSERT INTO test_model (name, created_at) VALUES ($1, $2) RETURNING *"
	if q.stmt != wantStmt {
		t.Errorf("q.stmt = %s, want %s", q.stmt, wantStmt)
	}

	wantParams := []interface{}{"Donkey Hote", createdAt}
	if len(q.params) != len(wantParams) {
		t.Errorf("len(q.params) = %d, want %d", len(q.params), len(wantParams))
		return
	}

	for i, param := range q.params {
		if param != wantParams[i] {
			t.Errorf("q.params[%d] = %v, want %v", i, param, wantParams[i])
		}
	}
}
