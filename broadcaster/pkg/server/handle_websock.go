package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

	"go.uber.org/zap"
)

func (s *Server) sub(cl *wsClient) {
	s.connsL.Lock()
	defer s.connsL.Unlock()
	s.conns[cl] = struct{}{}
}

func (s *Server) unsub(cl *wsClient) {
	s.connsL.Lock()
	defer s.connsL.Unlock()
	delete(s.conns, cl)
}

var errSendChannelFull = errors.New("send channel full")

func (cl *wsClient) sendBytes(
	ctx context.Context,
	body []byte,
) error {
	select {
	case cl.send <- body:
		return nil
	default:
		return errSendChannelFull
	}
}

func (cl *wsClient) sendMessage(
	ctx context.Context,
	ty string,
	payload interface{},
) error {
	body, err := json.Marshal(&wsMessage{
		Type:    ty,
		Payload: payload,
	})
	if err != nil {
		return err
	}
	return cl.sendBytes(ctx, body)
}

func (s *Server) getSubs(lock bool) []*wsClient {
	if lock {
		s.connsL.Lock()
		defer s.connsL.Unlock()
	}
	conns := make([]*wsClient, 0, len(s.conns))
	for c := range s.conns {
		conns = append(conns, c)
	}
	return conns
}

func (s *Server) broadcastLocal(
	ctx context.Context,
	body []byte,
) {
	s.connsL.Lock()
	defer s.connsL.Unlock()
	subs := s.getSubs(false)
	wg := new(sync.WaitGroup)
	wg.Add(len(subs))
	defer wg.Wait()
	for _, cl := range subs {
		func(cl *wsClient) {
			s.spawn(func() {
				defer wg.Done()
				if err := cl.sendBytes(
					ctx,
					body,
				); err != nil && err != errSendChannelFull {
					_ = cl.c.Close()
				}
			})
		}(cl)
	}
}

type wsMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload,omitempty"`
}

type wsClient struct {
	connID string
	ctx    context.Context
	c      *websocket.Conn
	log    *zap.Logger
	send   chan []byte
}

func (cl *wsClient) writePump() {
	for {
		select {
		case <-cl.ctx.Done():
			return
		case msg, ok := <-cl.send:
			if !ok {
				return
			}
			if err := cl.c.WriteMessage(
				websocket.TextMessage,
				msg,
			); err != nil {
				return
			}
		}
	}
}

func (s *Server) handleWebSock() http.HandlerFunc {
	return s.rbacHandler(
		http.MethodGet,
		nil,
		func(userID string, w http.ResponseWriter, r *http.Request) error {
			s.wg.Add(1)
			defer s.wg.Done()
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			upgrader := websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					// TODO: fix CORS check
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
			send := make(chan []byte, 256)
			cl := &wsClient{
				connID: connID,
				ctx:    r.Context(),
				c:      c,
				send:   send,
				log:    reqLog,
			}
			s.spawn(cl.writePump)
			defer c.Close()
			s.sub(cl)
			defer s.unsub(cl)
			if err := cl.sendMessage(
				r.Context(),
				"ping",
				nil,
			); err != nil {
				return fmt.Errorf("ping: %v", err)
			}
			for {
				select {
				case <-r.Context().Done():
					return r.Context().Err()
				case <-time.After(10 * time.Second):
					if err := cl.sendMessage(
						r.Context(),
						"ping",
						nil,
					); err != nil {
						return fmt.Errorf("ping: %v", err)
					}
				}
			}
		},
	)
}
