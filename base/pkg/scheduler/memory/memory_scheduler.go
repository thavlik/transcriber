package memory

import (
	"sync"

	"github.com/thavlik/transcriber/base/pkg/scheduler"
)

type memoryScheduler struct {
	set    map[string]struct{}
	locks  map[string]*memoryLock
	l      sync.Mutex
	notify []chan<- struct{}
}

func NewMemoryScheduler() scheduler.Scheduler {
	return &memoryScheduler{
		set:   make(map[string]struct{}),
		locks: make(map[string]*memoryLock),
	}
}
