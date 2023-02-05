package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/base/pkg/iam"
	pharmaseer "github.com/thavlik/transcriber/pharmaseer/pkg/api"
	"go.uber.org/zap"
)

func (s *Server) handleDrug() http.HandlerFunc {
	return s.rbacHandler(
		http.MethodGet,
		iam.NullPermissions,
		func(userID string, w http.ResponseWriter, r *http.Request) error {
			query := r.URL.Query().Get("q")
			if query == "" {
				return errors.New("missing query parameter 'q'")
			}
			s.log.Debug("querying drug details", zap.String("query", query))
			start := time.Now()
			drug, err := s.pharmaSeer.GetDrugDetails(
				r.Context(),
				pharmaseer.GetDrugDetails{
					Query: query,
				},
			)
			if err != nil {
				return errors.Wrap(err, "pharmaseer")
			}
			s.log.Debug("queried drug details",
				base.Elapsed(start),
				zap.Int("numSynonyms", len(drug.Synonyms)))
			w.Header().Set("Content-Type", "application/json")
			return json.NewEncoder(w).Encode(drug)
		},
	)
}

func (s *Server) handleDrugSvg() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			w.Header().Set("Access-Control-Allow-Origin", s.corsHeader)
			if r.Method == http.MethodOptions {
				base.AddPreflightHeaders(w)
				return nil
			} else if r.Method != http.MethodGet {
				return errors.Errorf("invalid method %s", r.Method)
			}
			id, ok := mux.Vars(r)["id"]
			if !ok {
				return errors.New("missing path parameter 'id'")
			}
			req, err := http.NewRequestWithContext(
				r.Context(),
				http.MethodGet,
				fmt.Sprintf("%s/structure?id=%s", s.pharmaSeerOpts.Endpoint, id),
				nil,
			)
			if err != nil {
				return err
			}
			resp, err := (&http.Client{
				Timeout: s.pharmaSeerOpts.Timeout,
			}).Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				return errors.Errorf(
					"failed to download svg for %s: status code %d",
					id,
					resp.StatusCode,
				)
			}
			w.Header().Set("Content-Type", "image/svg+xml")
			if _, err := io.Copy(w, resp.Body); err != nil {
				return err
			}
			return nil
		}(); err != nil {
			s.log.Error(r.RequestURI, zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (s *Server) getPDB(
	ctx context.Context,
	id string,
) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/pdb?id=%s", s.pharmaSeerOpts.Endpoint, id),
		nil,
	)
	if err != nil {
		return nil, err
	}
	resp, err := (&http.Client{
		Timeout: s.pharmaSeerOpts.Timeout,
	}).Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf(
			"failed to download pdb for %s: status code %d",
			id,
			resp.StatusCode,
		)
	}
	return resp.Body, nil
}

func (s *Server) handleDrugPdb() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			w.Header().Set("Access-Control-Allow-Origin", s.corsHeader)
			if r.Method == http.MethodOptions {
				base.AddPreflightHeaders(w)
				return nil
			} else if r.Method != http.MethodGet {
				return errors.Errorf("invalid method %s", r.Method)
			}
			id, ok := mux.Vars(r)["id"]
			if !ok {
				return errors.New("missing path parameter 'id'")
			}
			pdb, err := s.getPDB(r.Context(), id)
			if err != nil {
				return err
			}
			defer pdb.Close()
			w.Header().Set("Content-Type", "chemical/x-pdb")
			if _, err := io.Copy(w, pdb); err != nil {
				return err
			}
			return nil
		}(); err != nil {
			s.log.Error(r.RequestURI, zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (s *Server) handleDrugStl() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			w.Header().Set("Access-Control-Allow-Origin", s.corsHeader)
			if r.Method == http.MethodOptions {
				base.AddPreflightHeaders(w)
				return nil
			} else if r.Method != http.MethodGet {
				return errors.Errorf("invalid method %s", r.Method)
			}
			id, ok := mux.Vars(r)["id"]
			if !ok {
				return errors.New("missing path parameter 'id'")
			}
			pdb, err := s.getPDB(r.Context(), id)
			if err != nil {
				return err
			}
			defer pdb.Close()
			req, err := http.NewRequestWithContext(
				r.Context(),
				http.MethodPost,
				fmt.Sprintf("%s/convert", s.pdbMesh.Endpoint),
				pdb,
			)
			if err != nil {
				return err
			}
			resp, err := (&http.Client{
				Timeout: s.pdbMesh.Timeout,
			}).Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				return errors.Errorf(
					"failed to download svg for %s: status code %d",
					id,
					resp.StatusCode,
				)
			}
			w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
			if _, err := io.Copy(w, resp.Body); err != nil {
				return err
			}
			return nil
		}(); err != nil {
			s.log.Error(r.RequestURI, zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
