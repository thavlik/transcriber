package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/imgsearch/pkg/imgsearch/adapter"

	"go.uber.org/zap"
)

func (s *Server) isDisease(
	ctx context.Context,
	query string,
) (bool, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		s.define.Endpoint+"/disease?q="+url.QueryEscape(query),
		nil,
	)
	if err != nil {
		return false, err
	}
	resp, err := (&http.Client{
		Timeout: s.define.Timeout,
	}).Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("define service status code %d: %s", resp.StatusCode, string(body))
	}
	var result struct {
		IsDisease bool `json:"isDisease"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}
	return result.IsDisease, nil
}

func (s *Server) filterDiseaseQueryExpansion(
	ctx context.Context,
	queryExpansion []string,
) (result []string, err error) {
	wg := new(sync.WaitGroup)
	wg.Add(len(queryExpansion))
	defer wg.Wait()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	dones := make([]chan interface{}, len(queryExpansion))
	for i, query := range queryExpansion {
		done := make(chan interface{}, 1)
		dones[i] = done
		go func(i int, query string, done chan<- interface{}) {
			defer wg.Done()
			isDisease, err := s.isDisease(ctx, query)
			if err != nil {
				done <- err
				return
			}
			done <- isDisease
		}(i, query, done)
	}
	for i, done := range dones {
		select {
		case <-ctx.Done():
			cancel()
			return nil, ctx.Err()
		case v := <-done:
			switch v := v.(type) {
			case error:
				cancel()
				return nil, errors.Errorf(
					"failed to check if '%s' is a disease: %v",
					queryExpansion[i],
					v,
				)
			case bool:
				if v {
					result = append(result, queryExpansion[i])
				}
			default:
				panic(base.Unreachable)
			}
		}
	}
	return result, nil
}

func (s *Server) handleSearch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		retCode := http.StatusInternalServerError
		if err := func() error {
			var req struct {
				Query  string `json:"query"`
				UserID string `json:"userID,omitempty"`
			}
			switch r.Method {
			case http.MethodOptions:
				base.AddPreflightHeaders(w)
				return nil
			case http.MethodGet:
				req.Query = r.URL.Query().Get("q")
			case http.MethodPost:
				if r.Header.Get("Content-Type") != "application/json" {
					retCode = http.StatusUnsupportedMediaType
					return fmt.Errorf("unsupported media type %s", r.Header.Get("Content-Type"))
				}
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					retCode = http.StatusBadRequest
					return errors.Wrap(err, "failed to decode request body")
				}
			default:
				retCode = http.StatusMethodNotAllowed
				return errors.New("bad method")
			}
			w.Header().Set("Access-Control-Allow-Origin", "*")
			if req.Query == "" {
				retCode = http.StatusBadRequest
				return fmt.Errorf("missing query")
			}
			ty := r.URL.Query().Get("t")
			start := time.Now()
			service := adapter.Bing
			result, err := adapter.Search(
				r.Context(),
				service,
				req.Query,
				s.endpoint,
				s.apiKey,
				10,
				0,
			)
			if err != nil {
				return errors.Wrap(err, "search failed")
			}
			s.log.Debug(
				"searched for images",
				base.Elapsed(start),
				zap.String("service", string(service)),
				zap.Int("count", len(result.Images)),
			)
			w.Header().Set("Content-Type", "application/json")
			switch ty {
			case "DX_NAME":
				// only list query expansion terms that are diseases
				result.QueryExpansions, err = s.filterDiseaseQueryExpansion(
					r.Context(),
					result.QueryExpansions,
				)
				if err != nil {
					return errors.Wrap(err, "failed to filter query expansion")
				}
			default:
				// don't give query expansions
				result.QueryExpansions = nil
			}
			if err := json.NewEncoder(w).Encode(result); err != nil {
				return err
			}
			if req.UserID != "" {
				s.spawn(func() {
					// Only count searches for logged in users,
					// don't block the response, and only count
					// it if everything is successful. This is
					// the most generous policy for the user.
					s.pushSearchHistory(req.Query, req.UserID)
				})
			}
			return nil
		}(); err != nil {
			s.log.Error(r.RequestURI, zap.Error(err))
			w.WriteHeader(retCode)
			w.Write([]byte(err.Error()))
		}
	}
}
