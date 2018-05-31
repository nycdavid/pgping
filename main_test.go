package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"
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
	Pinged   int
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

type MockDelayer struct {
	i         int
	DelayFunc func()
}

func (m *MockDelayer) Delay() {
	m.DelayFunc()
}

func TestMain_exitNonZeroWithoutEnvVar(t *testing.T) {
	os.Setenv("PINGLIMIT", "1")
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
	dly := &MockDelayer{}
	dly.DelayFunc = func() {
		dly.i++
	}
	got := realMain(&MockLogger{buf: &buf}, c, dly)

	if expected != got {
		msg := fmt.Sprintf("Expected exit code %d, got %d.", expected, got)
		t.Error(msg)
	}
	if buf.String() == "" {
		t.Error("Expected buffer to have error")
	}
}

func TestMain_successfullyPingsTheDatabase(t *testing.T) {
	os.Setenv("PINGLIMIT", "1")
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
	dly := &MockDelayer{}
	dly.DelayFunc = func() {
		dly.i++
	}
	status := realMain(&MockLogger{buf: &buf}, c, dly)

	if status != 0 {
		t.Error("Expected a zero status code")
	}
	if buf.String() == "" {
		t.Error("Expected success log")
	}
	os.Unsetenv("PGCONN")
}

func TestMain_unsuccessfullyPingsTheDatabase(t *testing.T) {
	os.Setenv("PINGLIMIT", "1")
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
	dly := &MockDelayer{}
	dly.DelayFunc = func() {
		dly.i++
	}
	status := realMain(&MockLogger{buf: &buf}, c, dly)

	if status != 1 {
		t.Error("Expected a non-zero status code")
	}
	os.Unsetenv("PGCONN")
}

func TestMain_writesToLoggerWhenErroring(t *testing.T) {
	os.Setenv("PINGLIMIT", "1")
	var buf bytes.Buffer
	ml := &MockLogger{buf: &buf}
	mc := &MockConnection{}
	dly := &MockDelayer{}
	dly.DelayFunc = func() {
		dly.i++
	}
	realMain(ml, mc, dly)

	if ml.buf.String() == "" {
		t.Error("Expected non-empty log buffer")
	}
}

func TestMain_returnsNonZeroWhenSqlOpenErrors(t *testing.T) {
	os.Setenv("PINGLIMIT", "1")
	var buf bytes.Buffer
	os.Setenv("PGCONN", "pgconnstring")
	ml := &MockLogger{buf: &buf}
	mc := &MockConnection{
		OpenFunc: func() (SqlDB, error) {
			return nil, errors.New("Error")
		},
	}
	dly := &MockDelayer{}
	dly.DelayFunc = func() {
		dly.i++
	}
	if realMain(ml, mc, dly) == 0 {
		t.Error("Expected bad connection string to get non-zero code")
	}
	if ml.buf.String() == "" {
		t.Error("Expected error to be logged")
	}
	os.Unsetenv("PGCONN")
}

func TestMain_repeatedlyPingsUpToLimit(t *testing.T) {
	pingLmt := "5"
	os.Setenv("PINGLIMIT", pingLmt)
	var buf bytes.Buffer
	os.Setenv("PGCONN", "pgconnstring")
	ml := &MockLogger{buf: &buf}
	db := &MockSqlDB{}

	db.PingFunc = func() error {
		db.Pinged++
		return errors.New("Keep pinging")
	}
	mc := &MockConnection{
		OpenFunc: func() (SqlDB, error) {
			return db, nil
		},
	}
	dly := &MockDelayer{}
	dly.DelayFunc = func() {
		dly.i++
	}

	realMain(ml, mc, dly)

	if pingLmt != strconv.Itoa(db.Pinged) {
		msg := fmt.Sprintf("Expected db.Pinged to be %s, got %d", pingLmt, db.Pinged)
		t.Error(msg)
	}
	if dly.i == 0 {
		msg := fmt.Sprintf("Expected DelayFunc to be called > 0 times, got %d", dly.i)
		t.Error(msg)
	}
	os.Unsetenv("PGCONN")
}
