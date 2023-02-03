package server

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/thavlik/transcriber/scribe/pkg/source"
	"github.com/yutopp/go-rtmp"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (s *Server) listenRTMP(
	ctx context.Context,
	port int,
) error {
	s.wg.Add(1)
	defer s.wg.Done()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	tcpAddr, err := net.ResolveTCPAddr(
		"tcp",
		fmt.Sprintf(":%d", port),
	)
	if err != nil {
		return errors.Wrap(err, "resolve failed")
	}

	listener, err := net.ListenTCP(
		"tcp",
		tcpAddr,
	)
	if err != nil {
		return errors.Wrap(err, "listen failed")
	}

	srv := rtmp.NewServer(&rtmp.ServerConfig{
		OnConnect: func(conn net.Conn) (io.ReadWriteCloser, *rtmp.ConnConfig) {
			newSource := make(chan source.Source, 1)
			h := NewHandler(
				ctx,
				newSource,
				func(key string) error {
					if key == s.streamKey {
						return nil
					}
					return errors.New("invalid stream key")
				},
				s.wg,
				s.log,
			)
			s.spawn(func() {
				select {
				case <-h.ctx.Done():
					return
				case audioSource := <-newSource:
					select {
					case <-h.ctx.Done():
						return
					case s.newSource <- audioSource:
						break
					}
				}
			})
			return conn, &rtmp.ConnConfig{
				Handler: h,
				ControlState: rtmp.StreamControlStateConfig{
					DefaultBandwidthWindowSize: 6 * 1024 * 1024 / 8,
				},
			}
		},
	})

	s.spawn(func() {
		<-ctx.Done()
		_ = srv.Close()
	})

	s.log.Info(
		"rtmp server listening forever",
		zap.Int("port", port),
	)

	if err := srv.Serve(listener); err != nil {
		cancel()
		return err
	}

	return nil
}
