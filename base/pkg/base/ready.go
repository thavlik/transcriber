package base

import (
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var (
	start     time.Time = time.Now()
	readyFile string    = "/etc/ready"
)

func SignalReady(log *zap.Logger) {
	log.Debug("signaling ready", Elapsed(start))
	if err := os.WriteFile(
		readyFile,
		[]byte{1},
		0644,
	); err != nil {
		panic(errors.Wrap(err, "failed to write ready file"))
	}
}

func ReadyHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := os.Stat(readyFile); err == nil {
		return
	}
	w.WriteHeader(http.StatusNotFound)
}
