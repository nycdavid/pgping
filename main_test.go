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
	h      DBHandle
}

func (mc *MockConnection) Open(driverName, dataSourceName string) (*sql.DB, error) {
	mc.Opened = true
	if dataSourceName == "badpgconnstring" {
		return nil, errors.New("Bad dataSourceName")
	}
	return &sql.DB{}, nil
}

type MockLogger struct {
	buf *bytes.Buffer
	log []string
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

type MockDBHandle struct{}

func (mdbh *MockDBHandle) Ping() error {
	return nil
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
	ml := &MockLogger{buf: &buf}
	realMain(mc, ml)

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

func TestMain_pingsPgConnection(t *testing.T) {
	var buf bytes.Buffer
	os.Setenv("PGCONN", "goodpgconnstring-badping")
	ml := &MockLogger{buf: &buf}
	mc := &MockConnection{}

	if realMain(mc, ml) == 0 {
		t.Error("Expected bad ping to get non-zero code")
	}

	os.Unsetenv("PGCONN")
}
