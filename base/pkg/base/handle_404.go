package base

import (
	"net/http"

	"go.uber.org/zap"
)

func Handle404(
	log *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		AddPreflightHeaders(w)
		msg := "404: the requested page could not be found"
		log.Error(r.RequestURI, zap.String("err", msg))
		http.Error(w, msg, http.StatusNotFound)
	}
}
