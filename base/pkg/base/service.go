package base

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type BasicAuth struct {
	Username string
	Password string
}

type ServiceOptions struct {
	Endpoint  string
	Timeout   time.Duration
	BasicAuth BasicAuth
}

func (o *ServiceOptions) HasBasicAuth() bool {
	return o.BasicAuth.Username != ""
}

func (o *ServiceOptions) HasTimeout() bool {
	return o.Timeout != 0
}

func (o *ServiceOptions) WithTimeout(timeout time.Duration) *ServiceOptions {
	v := *o
	v.Timeout = timeout
	return &v
}

func (o *ServiceOptions) WithBasicAuth(username, password string) *ServiceOptions {
	v := *o
	v.BasicAuth = BasicAuth{
		Username: username,
		Password: password,
	}
	return &v
}

func AddServiceFlags(
	cmd *cobra.Command,
	name string,
	out *ServiceOptions,
	defaultTimeout time.Duration,
) {
	var defaultEndpoint string
	desc := name
	if name != "" {
		desc += " "
		name += "-"
	} else {
		// services with commands are expected to be
		// executed from within their pod, where the
		// server is listening on :80
		defaultEndpoint = "http://localhost:80"
	}
	cmd.PersistentFlags().StringVar(&out.Endpoint, name+"endpoint", defaultEndpoint, desc+"service endpoint")
	cmd.PersistentFlags().DurationVar(&out.Timeout, name+"timeout", defaultTimeout, desc+"service request timeout")
	cmd.PersistentFlags().StringVar(&out.BasicAuth.Username, name+"username", "", desc+"service basic auth username")
	cmd.PersistentFlags().StringVar(&out.BasicAuth.Password, name+"password", "", desc+"service basic auth password")
}

func ServiceEnv(name string, o *ServiceOptions) *ServiceOptions {
	if o == nil {
		o = &ServiceOptions{}
	}
	upper := strings.ToUpper(name)
	if len(upper) != 0 {
		upper = strings.ReplaceAll(upper, "-", "_") + "_"
	}
	CheckEnv(upper+"ENDPOINT", &o.Endpoint)
	if o.Endpoint == "" {
		if name != "" {
			panic(fmt.Sprintf("missing --%s-endpoint", name))
		} else {
			panic("missing --endpoint")
		}
	}
	timeoutKey := upper + "TIMEOUT"
	if v, ok := os.LookupEnv(timeoutKey); ok {
		timeout, err := time.ParseDuration(v)
		if err != nil {
			panic(fmt.Sprintf("parse %s from env: %v", timeoutKey, err))
		}
		o.Timeout = timeout
	}
	CheckEnv(upper+"USERNAME", &o.BasicAuth.Username)
	CheckEnv(upper+"PASSWORD", &o.BasicAuth.Password)
	return o
}
