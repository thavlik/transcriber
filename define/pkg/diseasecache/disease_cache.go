package diseasecache

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("no record found for input")

type DiseaseCache interface {
	Set(ctx context.Context, input string, isDisease bool) error
	IsDisease(ctx context.Context, input string) (bool, error)
}
