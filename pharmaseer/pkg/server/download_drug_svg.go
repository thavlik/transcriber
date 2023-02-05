package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/pharmaseer/pkg/thumbcache"
)

func downloadDrugSVG(
	ctx context.Context,
	drugBankAccessionNumber string,
	svgCache thumbcache.ThumbCache,
	w io.Writer,
) error {
	if has, err := svgCache.Has(
		ctx,
		drugBankAccessionNumber,
	); err != nil {
		return err
	} else if has {
		return nil
	}
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf(
			"https://go.drugbank.com/structures/%s/thumb.svg",
			drugBankAccessionNumber,
		),
		nil,
	)
	if err != nil {
		return err
	}
	resp, err := (&http.Client{
		Timeout: 10 * time.Second,
	}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf(
			"failed to download svg for %s: %s",
			drugBankAccessionNumber,
			resp.Status,
		)
	}
	var r io.Reader = resp.Body
	if w != nil {
		r = io.TeeReader(r, w)
	}
	if err := svgCache.Set(
		drugBankAccessionNumber,
		r,
	); err != nil {
		return err
	}
	return nil
}
