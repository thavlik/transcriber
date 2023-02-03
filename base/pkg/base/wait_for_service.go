package base

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

var maxServiceWaitTime = 30 * time.Second

func checkService(
	ctx context.Context,
	opts *ServiceOptions,
) error {
	url := fmt.Sprintf("%s/readyz", opts.Endpoint)
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return err
	}
	resp, err := (&http.Client{
		Timeout: 3 * time.Second,
	}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return fmt.Errorf("status code %d", resp.StatusCode)
}

func WaitForService(
	ctx context.Context,
	opts *ServiceOptions,
) {
	start := time.Now()
	var err error
	for time.Since(start) < maxServiceWaitTime {
		if err = checkService(ctx, opts); err == nil {
			return
		}
		time.Sleep(3 * time.Second)
	}
	panic(fmt.Errorf(
		"service %s not ready after %s, last error: %v",
		opts.Endpoint,
		maxServiceWaitTime.String(),
		err,
	))
}

func WaitForServices(
	ctx context.Context,
	opts []*ServiceOptions,
) {
	wg := new(sync.WaitGroup)
	wg.Add(len(opts))
	defer wg.Done()
	for _, o := range opts {
		go func(o *ServiceOptions) {
			WaitForService(ctx, o)
		}(o)
	}
}
