package infocache

import (
	"context"
	"errors"

	"github.com/thavlik/transcriber/pharmaseer/pkg/api"
)

var (
	ErrCacheUnavailable = errors.New("cache unavailable")
)

type InfoCache interface {
	HasDrug(ctx context.Context, query string) (bool, error)
	GetDrug(ctx context.Context, query string) (*api.DrugDetails, error)
	GetDrugByDrugBankAccessionNumber(ctx context.Context, drugBankAccessionNumber string) (*api.DrugDetails, error)
	SetDrug(ctx context.Context, query string, details *api.DrugDetails) error
}
