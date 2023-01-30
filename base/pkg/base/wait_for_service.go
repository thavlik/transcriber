package base

import (
	"fmt"
	"net/http"
	"time"
)

var timeout = time.Minute

func checkService(opts *ServiceOptions) error {
	url := fmt.Sprintf("%s/readyz", opts.Endpoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	if opts.HasBasicAuth() {
		req.SetBasicAuth(
			opts.BasicAuth.Username,
			opts.BasicAuth.Password,
		)
	}
	resp, err := (&http.Client{
		Timeout: 4 * time.Second,
	}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return fmt.Errorf("status code %d", resp.StatusCode)
}

func WaitForService(opts *ServiceOptions) {
	start := time.Now()
	var err error
	for time.Since(start) < timeout {
		if err = checkService(opts); err == nil {
			return
		}
		time.Sleep(3 * time.Second)
	}
	panic(fmt.Errorf(
		"service %s not ready after %s, last error: %v",
		opts.Endpoint,
		timeout.String(),
		err,
	))
}

func WaitForServices(opts []*ServiceOptions) {
	dones := make([]<-chan struct{}, len(opts))
	for i, o := range opts {
		done := make(chan struct{}, 1)
		dones[i] = done
		go func(o *ServiceOptions, done chan<- struct{}) {
			WaitForService(o)
			done <- struct{}{}
		}(o, done)
	}
	for _, done := range dones {
		<-done
	}
}
