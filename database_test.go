package kin

import (
	"database/sql"
	"os"
	"strings"
	"testing"

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
