package api

import (
	"encoding/json"

	"github.com/thavlik/transcriber/base/pkg/base"
)

func NewPharmaSeerClientFromOptions(opts base.ServiceOptions) PharmaSeer {
	options := NewPharmaSeerClientOptions().SetTimeout(opts.Timeout)
	if opts.BasicAuth.Username != "" {
		options.SetBasicAuth(opts.BasicAuth.Username, opts.BasicAuth.Password)
	}
	return NewPharmaSeerClient(opts.Endpoint, options)
}

func (d *DrugDetails) AsMap() map[string]interface{} {
	body, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}
	doc := make(map[string]interface{})
	if err := json.Unmarshal(body, &doc); err != nil {
		panic(err)
	}
	return doc
}

func ConvertDrugDetails(input map[string]interface{}) *DrugDetails {
	body, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}
	details := new(DrugDetails)
	if err := json.Unmarshal(body, details); err != nil {
		panic(err)
	}
	return details
}
