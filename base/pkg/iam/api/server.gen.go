// Code generated by oto; DO NOT EDIT.

package api

import (
	"context"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/pacedotdev/oto/otohttp"
)

var (
	remoteIAMDeleteUserTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "remote_i_a_m_delete_user_total",
		Help: "Auto-generated metric incremented on every call to RemoteIAM.DeleteUser",
	})
	remoteIAMDeleteUserSuccessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "remote_i_a_m_delete_user_success_total",
		Help: "Auto-generated metric incremented on every call to RemoteIAM.DeleteUser that does not return with an error",
	})

	remoteIAMLoginTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "remote_i_a_m_login_total",
		Help: "Auto-generated metric incremented on every call to RemoteIAM.Login",
	})
	remoteIAMLoginSuccessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "remote_i_a_m_login_success_total",
		Help: "Auto-generated metric incremented on every call to RemoteIAM.Login that does not return with an error",
	})
)

type RemoteIAM interface {
	DeleteUser(context.Context, DeleteUser) (*Void, error)
	Login(context.Context, LoginRequest) (*LoginResponse, error)
}

type remoteIAMServer struct {
	server    *otohttp.Server
	remoteIAM RemoteIAM
}

func RegisterRemoteIAM(server *otohttp.Server, remoteIAM RemoteIAM) {
	handler := &remoteIAMServer{
		server:    server,
		remoteIAM: remoteIAM,
	}
	server.Register("RemoteIAM", "DeleteUser", handler.handleDeleteUser)
	server.Register("RemoteIAM", "Login", handler.handleLogin)
}

func (s *remoteIAMServer) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	remoteIAMDeleteUserTotal.Inc()
	var request DeleteUser
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.remoteIAM.DeleteUser(r.Context(), request)
	if err != nil {
		log.Println("TODO: oto service error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	remoteIAMDeleteUserSuccessTotal.Inc()
}

func (s *remoteIAMServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	remoteIAMLoginTotal.Inc()
	var request LoginRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.remoteIAM.Login(r.Context(), request)
	if err != nil {
		log.Println("TODO: oto service error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	remoteIAMLoginSuccessTotal.Inc()
}

type DeleteUser struct {
	ID             string `json:"id"`
	Username       string `json:"username"`
	DeleteProjects bool   `json:"deleteProjects"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
	Error       string `json:"error,omitempty"`
}

type Void struct {
	Error string `json:"error,omitempty"`
}