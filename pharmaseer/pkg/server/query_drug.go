package server

import (
	"context"

	"github.com/thavlik/transcriber/pharmaseer/pkg/api"
)

const queryDrugScriptPath = "/scripts/query-drug.js"

func queryDrug(
	ctx context.Context,
	drugBankURL string,
	dest *api.DrugDetails,
) error {
	return nodeQuery(
		ctx,
		queryDrugScriptPath,
		drugBankURL,
		dest,
	)
}
