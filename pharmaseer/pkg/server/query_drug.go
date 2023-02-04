package server

import (
	"fmt"

	"github.com/thavlik/transcriber/pharmaseer/pkg/api"
)

const queryDrugScriptPath = "/scripts/query-drug.js"

func queryDrug(channelID string, dest *api.DrugDetails) error {
	return nodeQuery(
		queryDrugScriptPath,
		fmt.Sprintf("https://youtube.com/%s", channelID),
		dest,
	)
}
