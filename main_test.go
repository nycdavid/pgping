package main

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"
)

type MockConnection struct {
	Opened bool
}

func (mc *MockConnection) Open(driverName, dataSourceName string) (*sql.DB, error) {
	mc.Opened = true
	if dataSourceName == "badpgconnstring" {
		return &sql.DB{}, errors.New("Bad dataSourceName")
	}
	return &sql.DB{}, nil
}

type MockLogger struct {
	buf *bytes.Buffer
}

func (ml *MockLogger) Print(v ...interface{}) {
	var err error
	for _, i := range v {
		switch v := i.(type) {
		case string:
			_, err = ml.buf.WriteString(v)
		case error:
			_, err = ml.buf.WriteString(v.Error())
		}
		if err != nil {
			fmt.Println(err)
		}
	}
}

func TestMain_exitNonZeroWithoutEnvVar(t *testing.T) {
	var buf bytes.Buffer
	expected := 1
	got := realMain(&MockConnection{}, &MockLogger{buf: &buf})

	if expected != got {
		msg := fmt.Sprintf("Expected exit code %d, got %d.", expected, got)
		t.Error(msg)
	}
}

func TestMain_opensAPostgresConnection(t *testing.T) {
	var buf bytes.Buffer
	os.Setenv("PGCONN", "foo")
	mc := &MockConnection{}
	realMain(mc, &MockLogger{buf: &buf})

	if !mc.Opened {
		t.Error("Expected main to open a Postgres connection")
	}
	os.Unsetenv("PGCONN")
}

func TestMain_writesToLoggerWhenErroring(t *testing.T) {
	var buf bytes.Buffer
	ml := &MockLogger{buf: &buf}
	mc := &MockConnection{}
	realMain(mc, ml)

	if ml.buf.String() == "" {
		t.Error("Expected non-empty log buffer")
	}
}

func TestMain_returnsNonZeroWhenSqlOpenErrors(t *testing.T) {
	var buf bytes.Buffer
	os.Setenv("PGCONN", "badpgconnstring")
	ml := &MockLogger{buf: &buf}
	mc := &MockConnection{}
	if realMain(mc, ml) == 0 {
		t.Error("Expected bad connection string to get non-zero code")
	}
	os.Unsetenv("PGCONN")
}

func TestMain_writesToLogWhenSqlOpenErrors(t *testing.T) {
	var buf bytes.Buffer
	os.Setenv("PGCONN", "badpgconnstring")
	ml := &MockLogger{buf: &buf}
	mc := &MockConnection{}
	realMain(mc, ml)

	if ml.buf.String() == "" {
		t.Error("Expected error to be logged")
	}
	os.Unsetenv("PGCONN")
}
