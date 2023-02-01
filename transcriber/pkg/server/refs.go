package server

import (
	"time"

	"github.com/thavlik/transcriber/transcriber/pkg/refmat"
)

var useTimeout = 2 * time.Minute

func (s *Server) clearUsedRefs() {
	s.usedRefsL.Lock()
	defer s.usedRefsL.Unlock()
	s.usedRefs = make(map[*refmat.ReferenceMaterial]time.Time)
}

func (s *Server) isRefUsed(ref *refmat.ReferenceMaterial) bool {
	s.usedRefsL.Lock()
	lastUsed, ok := s.usedRefs[ref]
	s.usedRefsL.Unlock()
	return ok && time.Since(lastUsed) < useTimeout
}

func (s *Server) areRefsUsed(refs []*refmat.ReferenceMaterial) bool {
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

func (s *Server) useRef(ref *refmat.ReferenceMaterial) {
	s.usedRefsL.Lock()
	defer s.usedRefsL.Unlock()
	s.usedRefs[ref] = time.Now()
}
