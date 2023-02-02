package api

import "github.com/thavlik/transcriber/base/pkg/base"

func NewGatewayClientFromOptions(opts base.ServiceOptions) Gateway {
	options := NewGatewayClientOptions().SetTimeout(opts.Timeout)
	if opts.BasicAuth.Username != "" {
		options.SetBasicAuth(opts.BasicAuth.Username, opts.BasicAuth.Password)
	}
	return NewGatewayClient(opts.Endpoint, options)
}
