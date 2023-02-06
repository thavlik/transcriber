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

func downloadDrugBankPDB(
	ctx context.Context,
	drugBankAccessionNumber string,
	url string,
	force bool,
	pdbCache pdbcache.PDBCache,
	w io.Writer,
	log *zap.Logger,
) error {
	if !force {
		if has, err := pdbCache.Has(
			ctx,
			drugBankAccessionNumber,
		); err != nil {
			return errors.Wrap(err, "cache.Has")
		} else if has {
			log.Debug("cache has pdb")
			return nil
		}
	}
	log.Debug("downloading pdb from drugbank",
		zap.Bool("force", force))
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
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
	r := io.Reader(resp.Body)
	if w != nil {
		r = io.TeeReader(resp.Body, w)
	}
	if err := pdbCache.Set(
		drugBankAccessionNumber,
		r,
	); err != nil {
		return errors.Wrap(err, "set cache")
	}
	log.Debug("cached pdb")
	return nil
}
