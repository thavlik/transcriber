package server

import (
	"context"

	"github.com/thavlik/transcriber/pharmaseer/pkg/api"
)

const queryDrugScriptPath = "/scripts/query-drug.js"

func queryDrug(
	ctx context.Context,
	drugBankURL string,
) (*api.DrugDetails, error) {
	drug := new(api.DrugDetails)
	if err := nodeQuery(
		ctx,
		queryDrugScriptPath,
		drugBankURL,
		drug,
	); err != nil {
		return nil, err
	}
	return drug, nil
}
