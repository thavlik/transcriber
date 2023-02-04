package memory

import (
	"github.com/thavlik/transcriber/base/pkg/scheduler"
)

type memoryLock struct {
	projectID string
	s         *memoryScheduler
}

func (l *memoryLock) Extend() error {
	return nil
}

func (l *memoryLock) Release() error {
	l.s.unlock(l.projectID)
	return nil
}

func (s *memoryScheduler) Lock(projectID string) (scheduler.Lock, error) {
	s.l.Lock()
	defer s.l.Unlock()
	if _, ok := s.locks[projectID]; ok {
		return nil, scheduler.ErrLocked
	}
	v := &memoryLock{projectID, s}
	s.locks[projectID] = v
	return v, nil
}
func (s *memoryScheduler) unlock(projectID string) {
	s.l.Lock()
	defer s.l.Unlock()
	delete(s.locks, projectID)
}
