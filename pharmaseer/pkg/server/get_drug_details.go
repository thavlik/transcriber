package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/pharmaseer/pkg/api"
	"github.com/thavlik/transcriber/pharmaseer/pkg/infocache"
	"go.uber.org/zap"
)

func (s *Server) GetDrugDetails(
	ctx context.Context,
	req api.GetDrugDetails,
) (*api.DrugDetails, error) {
	log := s.log.With(zap.String("req.Input", req.Input))
	if req.Input == "" {
		return nil, errors.New("missing input")
	}
	cached, err := s.infoCache.GetDrug(ctx, req.Input)
	if req.Force || err == infocache.ErrCacheUnavailable {
		if err := s.scheduleDrugQuery(req.Input); err != nil {
			return nil, err
		}
	}
	if err == nil {
		log.Debug("drug details were cached")
		return cached, nil
	}
	return nil, errors.Wrap(err, "infocache.GetChannel")
}
