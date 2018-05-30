package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"
)

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

type MockSqlDB struct {
	PingFunc func() error
}

func (m *MockSqlDB) Ping() error {
	return m.PingFunc()
}

type MockConnection struct {
	OpenFunc func() (SqlDB, error)
}

func (m *MockConnection) Open(driverName, dataSourceName string) (SqlDB, error) {
	return m.OpenFunc()
}

func TestMain_exitNonZeroWithoutEnvVar(t *testing.T) {
	var buf bytes.Buffer
	expected := 1
	db := &MockSqlDB{
		PingFunc: func() error {
			return nil
		},
	}
	c := &MockConnection{
		OpenFunc: func() (SqlDB, error) {
			return db, nil
		},
	}
	got := realMain(&MockLogger{buf: &buf}, c)

	if expected != got {
		msg := fmt.Sprintf("Expected exit code %d, got %d.", expected, got)
		t.Error(msg)
	}
	if buf.String() == "" {
		t.Error("Expected buffer to have error")
	}
}

func TestMain_successfullyPingsTheDatabase(t *testing.T) {
	os.Setenv("PGCONN", "goodconnstring")
	var buf bytes.Buffer
	db := &MockSqlDB{
		PingFunc: func() error {
			return nil
		},
	}
	c := &MockConnection{
		OpenFunc: func() (SqlDB, error) {
			return db, nil
		},
	}
	status := realMain(&MockLogger{buf: &buf}, c)

	if status != 0 {
		t.Error("Expected a zero status code")
	}
	os.Unsetenv("PGCONN")
}

func TestMain_unsuccessfullyPingsTheDatabase(t *testing.T) {
	os.Setenv("PGCONN", "badconnstring")
	var buf bytes.Buffer
	db := &MockSqlDB{
		PingFunc: func() error {
			return errors.New("Bad ping")
		},
	}
	c := &MockConnection{
		OpenFunc: func() (SqlDB, error) {
			return db, nil
		},
	}
	status := realMain(&MockLogger{buf: &buf}, c)

	if status != 1 {
		t.Error("Expected a non-zero status code")
	}
	os.Unsetenv("PGCONN")
}

func TestMain_writesToLoggerWhenErroring(t *testing.T) {
	var buf bytes.Buffer
	ml := &MockLogger{buf: &buf}
	mc := &MockConnection{}
	realMain(ml, mc)

	if ml.buf.String() == "" {
		t.Error("Expected non-empty log buffer")
	}
}

func TestMain_returnsNonZeroWhenSqlOpenErrors(t *testing.T) {
	var buf bytes.Buffer
	os.Setenv("PGCONN", "pgconnstring")
	ml := &MockLogger{buf: &buf}
	mc := &MockConnection{
		OpenFunc: func() (SqlDB, error) {
			return nil, errors.New("Error")
		},
	}
	if realMain(ml, mc) == 0 {
		t.Error("Expected bad connection string to get non-zero code")
	}
	if ml.buf.String() == "" {
		t.Error("Expected error to be logged")
	}
	os.Unsetenv("PGCONN")
}
