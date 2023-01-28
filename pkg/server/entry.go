package server

import (
	"context"

	"github.com/gorilla/websocket"
	"github.com/thavlik/transcriber/pkg/refmat"
	"github.com/thavlik/transcriber/pkg/source"

	"go.uber.org/zap"
)

func Entry(
	httpPort int,
	rtmpPort int,
	metricsPort int,
	streamKey string,
	log *zap.Logger,
) error {
	go runMetrics(metricsPort, log)
	s := &server{
		newSource: make(chan source.Source, 16),
		l:         make(chan struct{}, 1),
		conns:     make(map[*websocket.Conn]struct{}),
		streamKey: streamKey,
		refs:      refmat.BuildReferenceMap(refmat.TestReferenceMaterials),
		usedRefs:  make(map[*refmat.ReferenceMaterial]struct{}),
		log:       log,
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case src := <-s.newSource:
				log.Info("new source")
				s.clearUsedRefs()
				if err := s.setSource(
					ctx,
					src,
				); err != nil {
					return
				}
			}
		}
	}()
	return s.ListenAndServe(
		httpPort,
		rtmpPort,
	)
}

func (s *server) clearUsedRefs() {
	s.usedRefsL.Lock()
	defer s.usedRefsL.Unlock()
	s.usedRefs = make(map[*refmat.ReferenceMaterial]struct{})
}

func (s *server) isRefUsed(ref *refmat.ReferenceMaterial) bool {
	s.usedRefsL.Lock()
	defer s.usedRefsL.Unlock()
	_, ok := s.usedRefs[ref]
	return ok
}

func (s *server) useRef(ref *refmat.ReferenceMaterial) {
	s.usedRefsL.Lock()
	defer s.usedRefsL.Unlock()
	s.usedRefs[ref] = struct{}{}
}
