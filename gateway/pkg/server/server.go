package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pacedotdev/oto/otohttp"
	remoteiam "github.com/thavlik/transcriber/base/pkg/iam/api"

	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/base/pkg/iam"

	"go.uber.org/zap"
)

type Server struct {
	iam        iam.IAM
	imgSearch  *base.ServiceOptions
	corsHeader string
	log        *zap.Logger
}

func NewServer(
	iam iam.IAM,
	imgSearch *base.ServiceOptions,
	corsHeader string,
	log *zap.Logger,
) *Server {
	s := &Server{
		iam,
		imgSearch,
		corsHeader,
		log,
	}
	return s
}

func (s *Server) AdminListenAndServe(port int) error {
	otoServer := otohttp.NewServer()
	remoteiam.RegisterRemoteIAM(otoServer, s)
	mux := http.NewServeMux()
	mux.Handle("/", otoServer)
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
	mux := http.NewServeMux()
	mux.HandleFunc("/", base.Handle404(s.log))
	mux.HandleFunc("/healthz", base.HealthHandler)
	mux.HandleFunc("/readyz", base.ReadyHandler)
	mux.HandleFunc("/user/login", s.handleLogin())
	mux.HandleFunc("/user/search", s.handleUserSearch())
	mux.HandleFunc("/user/signout", s.handleSignOut())
	mux.HandleFunc("/user/register", s.handleRegister())
	mux.HandleFunc("/user/resetpassword", s.handleSetPassword())
	mux.HandleFunc("/user/exists", s.handleUserExists())
	mux.HandleFunc("/img", s.handleImage())
	mux.HandleFunc("/img/search", s.handleImageSearch())
	s.log.Info("public api listening forever", zap.Int("port", port))
	return (&http.Server{
		Handler:      mux,
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
