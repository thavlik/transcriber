package storage

import (
	"context"
	"errors"
	"time"
)

type Definition struct {
	ID        string
	Input     string
	Output    string
	Timestamp time.Time
}

var ErrNotCached = errors.New("not cached")

type Storage interface {
	Insert(ctx context.Context, def *Definition) error
	GetMostRecent(ctx context.Context, input string) (*Definition, error)
}
