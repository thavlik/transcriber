package server

import (
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/pharmaseer/pkg/api"
	"github.com/thavlik/transcriber/pharmaseer/pkg/pdbcache"
	"go.uber.org/zap"
)

func (s *Server) handlePDB() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		retCode := http.StatusInternalServerError
		if err := func() error {
			if r.Method != http.MethodGet {
				retCode = http.StatusMethodNotAllowed
				return errors.New("method not allowed")
			}
			id := r.URL.Query().Get("id")
			if id == "" {
				retCode = http.StatusBadRequest
				return errors.New("missing query parameter 'id'")
			}
			reqLog := s.log.With(zap.String("id", id))
			if len(id) != 7 || !strings.HasPrefix(id, "DB") {
				reqLog.Warn("requested pdb for invalid drugbank")
				retCode = http.StatusBadRequest
				return errors.New("invalid drugbank accession number")
			}
			reqLog.Debug("retrieving pdb")
			drug, err := s.GetDrugDetails(
				r.Context(),
				api.GetDrugDetails{
					Query: id,
				},
			)
			if err != nil {
				return err
			}
			reqLog.Debug("resolved drugbank accession number",
				zap.String("chemicalFormula", drug.ChemicalFormula))
			w.Header().Set("Content-Type", "chemical/x-pdb")
			if err := s.pdbCache.Get(
				r.Context(),
				id,
				w,
			); err == pdbcache.ErrNotCached {
				return downloadPDB(
					r.Context(),
					drug.DrugBankAccessionNumber,
					drug.Structure.PDB,
					true,
					s.pdbCache,
					w,
					s.log,
				)
			} else if err != nil {
				return errors.Wrap(err, "pdbcache")
			}
			reqLog.Debug("downloaded pdb from cache")
			return nil
		}(); err != nil {
			s.log.Error(r.RequestURI, zap.Error(err))
			http.Error(w, err.Error(), retCode)
		}
	}
}
