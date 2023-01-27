package server

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/thavlik/transcriber/pkg/source/aac"
	"github.com/yutopp/go-rtmp"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (s *server) listenRTMP(port int) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return errors.Wrap(err, "resolve failed")
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return errors.Wrap(err, "listen failed")
	}
	srv := rtmp.NewServer(&rtmp.ServerConfig{
		OnConnect: func(conn net.Conn) (io.ReadWriteCloser, *rtmp.ConnConfig) {
			processAudio := make(chan []byte, 16)
			newSource := make(chan *aac.AACSource, 1)
			h := NewHandler(
				context.Background(),
				newSource,
				processAudio,
				func(key string) bool {
					return key == s.streamKey
				},
				s.log,
			)
			go func() {
				var audioSource *aac.AACSource
				select {
				case <-h.ctx.Done():
					return
				case audioSource = <-newSource:
					select {
					case <-h.ctx.Done():
						return
					case s.newSource <- audioSource:
						break
					}
				}
				for {
					select {
					case <-h.ctx.Done():
						return
					case audio, ok := <-processAudio:
						if !ok {
							return
						}
						if n, err := audioSource.Write(audio); err != nil {
							s.log.Error("failed to write audio bytes", zap.Error(err))
						} else if n != len(audio) {
							s.log.Error("failed to write all audio bytes", zap.Int("n", n), zap.Int("len", len(audio)))
						}
					}
				}
			}()
			return conn, &rtmp.ConnConfig{
				Handler: h,
				ControlState: rtmp.StreamControlStateConfig{
					DefaultBandwidthWindowSize: 6 * 1024 * 1024 / 8,
				},
			}
		},
	})
	s.log.Info("rtmp server listening forever", zap.Int("port", port))
	if err := srv.Serve(listener); err != nil {
		return errors.Wrap(err, "failed to serve")
	}
	return nil
}
