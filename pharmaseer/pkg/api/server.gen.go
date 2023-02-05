// Code generated by oto; DO NOT EDIT.

package api

import (
	"context"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/pacedotdev/oto/otohttp"
)

var (
	pharmaSeerGetDrugDetailsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pharma_seer_get_drug_details_total",
		Help: "Auto-generated metric incremented on every call to PharmaSeer.GetDrugDetails",
	})
	pharmaSeerGetDrugDetailsSuccessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pharma_seer_get_drug_details_success_total",
		Help: "Auto-generated metric incremented on every call to PharmaSeer.GetDrugDetails that does not return with an error",
	})
)

type PharmaSeer interface {
	GetDrugDetails(context.Context, GetDrugDetails) (*DrugDetails, error)
}

type pharmaSeerServer struct {
	server     *otohttp.Server
	pharmaSeer PharmaSeer
}

func RegisterPharmaSeer(server *otohttp.Server, pharmaSeer PharmaSeer) {
	handler := &pharmaSeerServer{
		server:     server,
		pharmaSeer: pharmaSeer,
	}
	server.Register("PharmaSeer", "GetDrugDetails", handler.handleGetDrugDetails)
}

func (s *pharmaSeerServer) handleGetDrugDetails(w http.ResponseWriter, r *http.Request) {
	pharmaSeerGetDrugDetailsTotal.Inc()
	var request GetDrugDetails
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.pharmaSeer.GetDrugDetails(r.Context(), request)
	if err != nil {
		log.Println("TODO: oto service error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	pharmaSeerGetDrugDetailsSuccessTotal.Inc()
}

type DrugDetails struct {
	Summary                 string            `json:"summary"`
	BrandNames              []string          `json:"brandNames"`
	GenericName             string            `json:"genericName"`
	DrugBankAccessionNumber string            `json:"drugBankAccessionNumber"`
	Background              string            `json:"background"`
	Groups                  []string          `json:"groups"`
	Structure               *DrugStructure    `json:"structure"`
	Weight                  *Weight           `json:"weight"`
	ChemicalFormula         string            `json:"chemicalFormula"`
	Synonyms                []string          `json:"synonyms"`
	ExternalIDs             []string          `json:"externalIDs"`
	Pharmacology            *DrugPharmacology `json:"pharmacology"`
	References              *References       `json:"references"`
	Error                   string            `json:"error,omitempty"`
}

type DrugPharmacology struct {
	Indication           string      `json:"indication"`
	AssociatedConditions []string    `json:"associatedConditions"`
	Pharmacodynamics     string      `json:"pharmacodynamics"`
	MechanismOfAction    string      `json:"mechanismOfAction"`
	Absorption           string      `json:"absorption"`
	VolumeOfDistribution string      `json:"volumeOfDistribution"`
	ProteinBinding       string      `json:"proteinBinding"`
	Metabolism           *Metabolism `json:"metabolism"`
	RouteOfElimination   string      `json:"routeOfElimination"`
	HalfLife             string      `json:"halfLife"`
	Clearance            string      `json:"clearance"`
	Toxicity             string      `json:"toxicity"`
}

type DrugStructure struct {
	Image string `json:"image"`
	PDB   string `json:"pDB"`
}

type ExperimentalMetric struct {
	Value  float64 `json:"value"`
	Source string  `json:"source"`
}

type ExperimentalProperties struct {
	MeltingPoint    ExperimentalMetric `json:"meltingPoint"`
	WaterSolubility ExperimentalMetric `json:"waterSolubility"`
	LogP            ExperimentalMetric `json:"logP"`
	Caco2Perm       ExperimentalMetric `json:"caco2Perm"`
	PKa             ExperimentalMetric `json:"pKa"`
}

type GetDrugDetails struct {
	Query string `json:"query"`
	Force bool   `json:"force"`
}

type Metabolism struct {
	Description string `json:"description"`
}

type Reference struct {
	Index int    `json:"index"`
	Title string `json:"title"`
	Link  string `json:"link"`
}

type References struct {
	General []*Reference `json:"general"`
}

type Weight struct {
	Average      string `json:"average"`
	Monoisotopic string `json:"monoisotopic"`
}
