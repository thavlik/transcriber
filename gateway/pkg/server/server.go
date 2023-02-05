package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/pacedotdev/oto/otohttp"
	remoteiam "github.com/thavlik/transcriber/base/pkg/iam/api"
	pharmaseer "github.com/thavlik/transcriber/pharmaseer/pkg/api"

	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/base/pkg/iam"

	"go.uber.org/zap"
)

type Server struct {
	ctx            context.Context
	cancel         context.CancelFunc
	iam            iam.IAM
	imgSearch      *base.ServiceOptions
	define         *base.ServiceOptions
	pharmaSeerOpts *base.ServiceOptions
	pharmaSeer     pharmaseer.PharmaSeer
	corsHeader     string
	wg             *sync.WaitGroup
	log            *zap.Logger
}

func NewServer(
	ctx context.Context,
	iam iam.IAM,
	imgSearch *base.ServiceOptions,
	define *base.ServiceOptions,
	pharmaSeerOpts *base.ServiceOptions,
	pharmaSeer pharmaseer.PharmaSeer,
	corsHeader string,
	log *zap.Logger,
) *Server {
	ctx, cancel := context.WithCancel(ctx)
	s := &Server{
		ctx,
		cancel,
		iam,
		imgSearch,
		define,
		pharmaSeerOpts,
		pharmaSeer,
		corsHeader,
		new(sync.WaitGroup),
		log,
	}
	return s
}

func (s *Server) spawn(f func()) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		f()
	}()
}

func (s *Server) ShutDown() {
	s.cancel()
	s.wg.Wait()
}

func (s *Server) AdminListenAndServe(port int) error {
	mux := http.NewServeMux()
	if s.iam != nil {
		otoServer := otohttp.NewServer()
		remoteiam.RegisterRemoteIAM(otoServer, s)
		mux.Handle("/", otoServer)
	}
	mux.HandleFunc("/healthz", base.HealthHandler)
	mux.HandleFunc("/readyz", base.ReadyHandler)
	s.log.Info("iam admin listening forever", zap.Int("port", port))
	return (&http.Server{
		Handler:      mux,
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}).ListenAndServe()
}

func (s *Server) ListenAndServe(port int) error {
	router := mux.NewRouter()
	router.HandleFunc("/", base.Handle404(s.log))
	router.HandleFunc("/healthz", base.Handle200)
	router.HandleFunc("/readyz", base.Handle200)
	router.HandleFunc("/define", s.handleDefine())
	router.HandleFunc("/disease", s.handleIsDisease())
	router.HandleFunc("/img", s.handleImage())
	router.HandleFunc("/img/search", s.handleImageSearch())
	if s.pharmaSeer != nil {
		router.HandleFunc("/drug", s.handleDrug())
		router.HandleFunc("/drug/{id}/structure.svg", s.handleDrugSvg())
	}
	if s.iam != nil {
		router.HandleFunc("/user/login", s.handleLogin())
		router.HandleFunc("/user/search", s.handleUserSearch())
		router.HandleFunc("/user/signout", s.handleSignOut())
		router.HandleFunc("/user/register", s.handleRegister())
		router.HandleFunc("/user/resetpassword", s.handleSetPassword())
		router.HandleFunc("/user/exists", s.handleUserExists())
	}
	s.log.Info("public api listening forever", zap.Int("port", port))
	return (&http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}).ListenAndServe()
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
			}
			if method != "" && r.Method != method {
				w.WriteHeader(http.StatusBadRequest)
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
