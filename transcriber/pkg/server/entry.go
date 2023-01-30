package server

import (
	"context"
	"time"

	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/transcriber/pkg/refmat"
	"github.com/thavlik/transcriber/transcriber/pkg/source"

	"go.uber.org/zap"
)

var useTimeout = 2 * time.Minute

func Entry(
	httpPort int,
	rtmpPort int,
	metricsPort int,
	streamKey string,
	log *zap.Logger,
) error {
	go base.RunMetrics(metricsPort, log)
	s := &server{
		newSource: make(chan source.Source, 16),
		l:         make(chan struct{}, 1),
		conns:     make(map[*wsClient]struct{}),
		streamKey: streamKey,
		refs:      refmat.BuildReferenceMap(refmat.TestReferenceMaterials),
		usedRefs:  make(map[*refmat.ReferenceMaterial]time.Time),
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
	s.usedRefs = make(map[*refmat.ReferenceMaterial]time.Time)
}

func (s *server) isRefUsed(ref *refmat.ReferenceMaterial) bool {
	s.usedRefsL.Lock()
	lastUsed, ok := s.usedRefs[ref]
	s.usedRefsL.Unlock()
	return ok && time.Since(lastUsed) < useTimeout
}

func (s *server) areRefsUsed(refs []*refmat.ReferenceMaterial) bool {
	s.usedRefsL.Lock()
	defer s.usedRefsL.Unlock()
	for _, ref := range refs {
		lastUsed, ok := s.usedRefs[ref]
		if !ok || time.Since(lastUsed) > useTimeout {
			return false
		}
	}
	return true
}

func (s *server) useRef(ref *refmat.ReferenceMaterial) {
	s.usedRefsL.Lock()
	defer s.usedRefsL.Unlock()
	s.usedRefs[ref] = time.Now()
}
