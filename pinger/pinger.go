package pinger

import (
	"errors"
)

type Pinger struct {
}

func NewPinger(cs string) (*Pinger, error) {
	if cs == "" {
		return &Pinger{}, errors.New("Expected non-empty connection string.")
	}
	return &Pinger{}, nil
}
