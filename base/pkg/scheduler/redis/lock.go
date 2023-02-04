package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/bsm/redislock"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/scheduler"
)

type redisLock struct {
	ttl  time.Duration
	lock *redislock.Lock
}

func (l *redisLock) Extend() error {
	if err := l.lock.Refresh(
		context.Background(),
		l.ttl,
		nil,
	); err != nil {
		return errors.Wrap(err, "redislock")
	}
	return nil
}

func (l *redisLock) Release() error {
	if err := l.lock.Release(context.Background()); err != nil {
		return errors.Wrap(err, "redislock")
	}
	return nil
}

func lockKey(entity string) string {
	return fmt.Sprintf("pl:{%s}", entity)
}

func (s *redisScheduler) Lock(entity string) (scheduler.Lock, error) {
	lock, err := s.locker.Obtain(
		context.Background(),
		lockKey(entity),
		s.lockTTL,
		nil,
	)
	if err == redislock.ErrNotObtained {
		return nil, scheduler.ErrLocked
	} else if err != nil {
		return nil, errors.Wrap(err, "redislock")
	}
	return &redisLock{
		lock: lock,
		ttl:  s.lockTTL,
	}, nil
}
