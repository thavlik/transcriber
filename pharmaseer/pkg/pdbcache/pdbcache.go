package pdbcache

import (
	"context"
	"errors"
	"io"
)

var ErrNotCached = errors.New("pdb not cached")

type PDBCache interface {
	Has(ctx context.Context, id string) (bool, error)
	Get(ctx context.Context, id string, w io.Writer) error
	Set(id string, r io.Reader) error
}
