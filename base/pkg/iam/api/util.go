package api

import "github.com/thavlik/transcriber/base/pkg/base"

func NewRemoteIAMClientFromOptions(opts base.ServiceOptions) RemoteIAM {
	options := NewRemoteIAMClientOptions().SetTimeout(opts.Timeout)
	if opts.BasicAuth.Username != "" {
		options.SetBasicAuth(opts.BasicAuth.Username, opts.BasicAuth.Password)
	}
	return NewRemoteIAMClient(opts.Endpoint, options)
}
