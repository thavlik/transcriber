package redis_disease_cache

import (
	"context"
	"sync"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

func (m *redisDiseaseCache) Set(
	ctx context.Context,
	input string,
	isDisease bool,
) error {
	var underlyingDone chan error
	if m.underlying != nil {
		wg := new(sync.WaitGroup)
		wg.Add(1)
		defer wg.Wait()
		underlyingDone = make(chan error, 1)
		go func() {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			case underlyingDone <- m.underlying.Set(
				ctx,
				input,
				isDisease,
			):
			}
		}()
	}
	var value string
	if isDisease {
		value = "1"
	} else {
		value = "0"
	}
	if _, err := m.redis.Set(
		ctx,
		key(input),
		value,
		0,
	).Result(); err != nil {
		if underlyingDone != nil {
			select {
			case <-ctx.Done():
				var multi error
				multi = multierror.Append(multi, err)
				multi = multierror.Append(multi, ctx.Err())
				return multi
			case err2 := <-underlyingDone:
				var multi error
				multi = multierror.Append(multi, err)
				multi = multierror.Append(multi, errors.Wrap(err2, "set underlying cache"))
				return multi
			}
		}
		return errors.Wrap(err, "redis")
	}
	if underlyingDone != nil {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-underlyingDone:
			return errors.Wrap(err, "set underlying cache")
		}
	}
	return nil
}
