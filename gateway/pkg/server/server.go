package server

import (
	"encoding/json"
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
	corsHeader string
	log        *zap.Logger
}

func NewServer(
	iam iam.IAM,
	corsHeader string,
	log *zap.Logger,
) *Server {
	s := &Server{
		iam,
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
	mux.HandleFunc("/healthz", base.HealthHandler)
	mux.HandleFunc("/readyz", base.ReadyHandler)
	mux.HandleFunc("/user/login", s.handleLogin())
	mux.HandleFunc("/user/search", s.handleUserSearch())
	mux.HandleFunc("/user/signout", s.handleSignOut())
	mux.HandleFunc("/user/register", s.handleRegister())
	mux.HandleFunc("/user/resetpassword", s.handleSetPassword())
	mux.HandleFunc("/user/exists", s.handleUserExists())
	s.log.Info("public api listening forever", zap.Int("port", port))
	return (&http.Server{
		Handler:      mux,
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}).ListenAndServe()
}

func addPreflightHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", "0")
	w.Header().Set("Access-Control-Max-Age", "1728000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "AccessToken,DNT,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type")
	w.WriteHeader(http.StatusNoContent)
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
				addPreflightHeaders(w)
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
			w.WriteHeader(http.StatusInternalServerError)
			s.log.Error("handler error",
				zap.String("r.RequestURI", r.RequestURI),
				zap.Error(err))
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
			addPreflightHeaders(w)
			return
		}
		if err := func() (err error) {
			if method != "" && r.Method != method {
				w.WriteHeader(http.StatusBadRequest)
				return nil
			}
			return f(w, r)
		}(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			s.log.Error("handler error",
				zap.Error(err),
				zap.String("r.RequestURI", r.RequestURI))
		}
	}
}

func writeError(
	w http.ResponseWriter,
	code int,
	err error,
) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
}
