package server

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/pharmaseer/pkg/pdbcache"
	"go.uber.org/zap"
)

func downloadPDB(
	ctx context.Context,
	pdb *pdbItem,
	pdbCache pdbcache.PDBCache,
	log *zap.Logger,
) error {
	if !pdb.Force {
		if has, err := pdbCache.Has(
			ctx,
			pdb.DrugBankAccessionNumber,
		); err != nil {
			return errors.Wrap(err, "cache.Has")
		} else if has {
			log.Debug("cache has pdb")
			return nil
		}
	}
	log.Debug("downloading pdb from drugbank",
		zap.Bool("force", pdb.Force))
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		pdb.URL,
		nil,
	)
	if err != nil {
		return err
	}
	resp, err := (&http.Client{
		Timeout: 20 * time.Second,
	}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return errors.Errorf(
			"unexpected status code %d: %s",
			resp.StatusCode,
			string(body),
		)
	}
	if err := pdbCache.Set(
		pdb.DrugBankAccessionNumber,
		resp.Body,
	); err != nil {
		return errors.Wrap(err, "set cache")
	}
	log.Debug("cached pdb")
	return nil
}
