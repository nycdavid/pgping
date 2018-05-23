package main

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
)

type MockConnection struct {
	Opened bool
}

func (mc *MockConnection) Open(driverName, dataSourceName string) (*sql.DB, error) {
	mc.Opened = true
	return &sql.DB{}, nil
}

func TestMain_exitNonZeroWithoutEnvVar(t *testing.T) {
	expected := 1
	got := realMain(&MockConnection{})

	if expected != got {
		msg := fmt.Sprintf("Expected exit code %d, got %d.", expected, got)
		t.Error(msg)
	}
}

func TestMain_opensAPostgresConnection(t *testing.T) {
	os.Setenv("PGCONN", "foo")
	mc := &MockConnection{}
	realMain(mc)

	if !mc.Opened {
		t.Error("Expected main to open a Postgres connection")
	}
}
