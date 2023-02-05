package server

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/pharmaseer/pkg/api"
	"github.com/thavlik/transcriber/pharmaseer/pkg/infocache"
)

func drugTopic(input string) string {
	return fmt.Sprintf("drug:%s", input)
}

func (s *Server) GetDrugDetails(
	ctx context.Context,
	req api.GetDrugDetails,
) (*api.DrugDetails, error) {
	if req.Query == "" {
		return nil, errors.New("missing query")
	}
	if !req.Force {
		// try and use the cache
		var cached *api.DrugDetails
		var err error
		if strings.HasPrefix(req.Query, "DB") && len(req.Query) == 7 {
			cached, err = s.infoCache.GetDrugByDrugBankAccessionNumber(
				ctx,
				req.Query,
			)
		} else {
			cached, err = s.infoCache.GetDrug(ctx, req.Query)
		}
		if err == nil {
			// drug was cached
			return cached, nil
		} else if err != infocache.ErrCacheUnavailable {
			return nil, err
		}
	}
	// subscribe to event where cache becomes available
	sub, err := s.pubsub.Subscribe(ctx, drugTopic(req.Query))
	if err != nil {
		return nil, err
	}
	defer sub.Cancel(s.ctx)
	msgs := sub.Messages(ctx)
	// schedule the drug query
	if err := s.scheduleDrugQuery(
		req.Query,
		req.Force,
	); err != nil {
		return nil, err
	}
	// wait for the query to complete
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case msg := <-msgs:
			details := new(api.DrugDetails)
			if err := json.Unmarshal(msg, details); err != nil {
				return nil, err
			}
			return details, nil
		}
	}
}
