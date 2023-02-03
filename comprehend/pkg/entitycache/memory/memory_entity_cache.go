package memory

import (
	"github.com/thavlik/transcriber/comprehend/pkg/comprehend"
	"github.com/thavlik/transcriber/comprehend/pkg/entitycache"
)

type memoryEntityCache struct {
	m map[string]*comprehend.Entity
	l chan struct{}
}

func NewMemoryEntityCache() entitycache.EntityCache {
	return &memoryEntityCache{
		make(map[string]*comprehend.Entity),
		make(chan struct{}, 1),
	}
}
