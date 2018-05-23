package pinger

import (
	"testing"
)

func TestExitsNonZeroWithoutConnectionString(t *testing.T) {
	_, err := NewPinger("")

	if err == nil {
		t.Error("Expected error with empty connection string.")
	}
}
