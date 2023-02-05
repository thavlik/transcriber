package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/thavlik/transcriber/base/pkg/iam"

	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/base/pkg/pubsub"

	"go.uber.org/zap"
)

type Server struct {
	ctx        context.Context
	cancel     context.CancelFunc
	corsHeader string
	conns      map[*wsClient]struct{}
	connsL     sync.Mutex
	pub        pubsub.Publisher
	wg         *sync.WaitGroup
	iam        iam.IAM
	log        *zap.Logger
}

func NewServer(
	ctx context.Context,
	iam iam.IAM,
	pub pubsub.Publisher,
	corsHeader string,
	log *zap.Logger,
) *Server {
	ctx, cancel := context.WithCancel(ctx)
	return &Server{
		ctx:        ctx,
		cancel:     cancel,
		iam:        iam,
		pub:        pub,
		corsHeader: corsHeader,
		wg:         new(sync.WaitGroup),
		conns:      make(map[*wsClient]struct{}),
		log:        log,
	}
}

func (s *Server) ListenAndServe(
	httpPort int,
) error {
	ctx, cancel := context.WithCancel(s.ctx)
	defer cancel()

	mux := http.NewServeMux()
	mux.HandleFunc("/", base.Handle404(s.log))
	mux.HandleFunc("/healthz", base.Handle200)
	mux.HandleFunc("/readyz", base.ReadyHandler)
	mux.HandleFunc("/ws", s.handleWebSock())
	mux.HandleFunc("/publish", s.handlePublish())

	srv := &http.Server{
		Handler:      mux,
		Addr:         fmt.Sprintf("0.0.0.0:%d", httpPort),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	s.spawn(func() {
		<-ctx.Done()
		_ = srv.Shutdown(ctx)
	})

	s.log.Info(
		"http server listening forever",
		zap.Int("port", httpPort),
	)
	return srv.ListenAndServe()
}

func (s *Server) ShutDown() {
	s.cancel()
	s.connsL.Lock()
	for conn := range s.conns {
		_ = conn.c.Close()
	}
	s.connsL.Unlock()
	s.wg.Wait()
}

func (s *Server) spawn(f func()) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		f()
	}()
}

func (s *Server) rbacHandler(
	method string,
	permissions []string,
	f func(userID string, w http.ResponseWriter, r *http.Request) error,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := func() (err error) {
			w.Header().Set("Access-Control-Allow-Origin", s.corsHeader)
			if r.Method == http.MethodOptions {
				base.AddPreflightHeaders(w)
				return nil
			} else if method != "" && r.Method != method {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return nil
			}
			var userID string
			if permissions != nil {
				// empty slice of permissions checks login
				// without requiring any specific permission
				userID, err = s.rbac(r.Context(), r, permissions)
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					s.log.Error("auth failure",
						zap.String("r.RequestURI", r.RequestURI),
						zap.Error(err))
					return nil
				}
			}
			return f(userID, w, r)
		}(); err != nil {
			s.log.Error(r.RequestURI, zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (s *Server) handler(
	method string,
	f func(w http.ResponseWriter, r *http.Request) error,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", s.corsHeader)
		if r.Method == http.MethodOptions {
			base.AddPreflightHeaders(w)
			return
		}
		if err := func() (err error) {
			if method != "" && r.Method != method {
				w.WriteHeader(http.StatusBadRequest)
				return nil
			}
			return f(w, r)
		}(); err != nil {
			s.log.Error(r.RequestURI, zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
