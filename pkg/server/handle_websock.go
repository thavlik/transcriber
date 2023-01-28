package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"go.uber.org/zap"
)

func addPreflightHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", "0")
	w.Header().Set("Access-Control-Max-Age", "1728000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "AccessToken,DNT,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type")
	w.WriteHeader(http.StatusNoContent)
}

func (s *server) sub(c *websocket.Conn) {
	s.connsL.Lock()
	defer s.connsL.Unlock()
	s.conns[c] = struct{}{}
}

func (s *server) unsub(c *websocket.Conn) {
	s.connsL.Lock()
	defer s.connsL.Unlock()
	delete(s.conns, c)
}

func (s *server) getSubs() []*websocket.Conn {
	s.connsL.Lock()
	defer s.connsL.Unlock()
	conns := make([]*websocket.Conn, 0, len(s.conns))
	for c := range s.conns {
		conns = append(conns, c)
	}
	return conns
}

func (s *server) broadcast(body []byte) {
	subs := s.getSubs()
	for _, sub := range subs {
		if err := sub.WriteMessage(
			websocket.TextMessage,
			body,
		); err != nil {
			s.log.Warn("failed to write message to websocket, closing connection", zap.Error(err))
			_ = sub.Close()
		}
	}
}

type wsMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func (s *server) handleWebSock() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		retCode := http.StatusInternalServerError
		if err := func() error {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			switch r.Method {
			case http.MethodOptions:
				addPreflightHeaders(w)
				return nil
			case http.MethodGet:
				break
			default:
				retCode = http.StatusMethodNotAllowed
				return fmt.Errorf("method not allowed")
			}
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			upgrader := websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			}
			c, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				return fmt.Errorf("upgrade: %v", err)
			}
			connID := uuid.New().String()
			reqLog := s.log.With(zap.String("connID", connID))
			reqLog.Debug("upgraded websocket connection")
			defer reqLog.Debug("closed websocket connection")
			defer c.Close()
			s.sub(c)
			defer s.unsub(c)
			ping, _ := json.Marshal(&wsMessage{Type: "ping"})
			if err := c.WriteMessage(
				websocket.TextMessage,
				ping,
			); err != nil {
				return fmt.Errorf("ping: %v", err)
			}
			for {
				select {
				case <-r.Context().Done():
					return r.Context().Err()
				case <-time.After(10 * time.Second):
					if err := c.WriteMessage(
						websocket.TextMessage,
						ping,
					); err != nil {
						return fmt.Errorf("ping: %v", err)
					}
				}
			}
		}(); err != nil {
			s.log.Error(r.RequestURI, zap.Error(err))
			w.WriteHeader(retCode)
			w.Write([]byte(err.Error()))
		}
	}
}
