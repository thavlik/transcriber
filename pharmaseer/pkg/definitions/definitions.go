package definitions

type PharmaSeer interface {
	GetDrugDetails(GetDrugDetails) DrugDetails
}

type GetDrugDetails struct {
	Input string `json:"input"`
	Force bool   `json:"force"`
}

type DrugStructure struct {
	ImageURL string `json:"imageURL"`
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

type DrugPharmacology struct {
	RouteOfElimination string `json:"routeOfElimination"`
	HalfLife           string `json:"halfLife"`
	MechanismOfAction  string `json:"mechanismOfAction"`
}

type DrugDetails struct {
	Summary                string                  `json:"summary"`
	BrandNames             []string                `json:"brandNames"`
	GenericName            string                  `json:"genericName"`
	Type                   string                  `json:"type"`
	Groups                 []string                `json:"groups"`
	AverageWeight          float64                 `json:"averageWeight"`
	MonoisotopicWeight     float64                 `json:"monoisotopicWeight"`
	ChemicalFormula        string                  `json:"chemicalFormula"`
	Structure              *DrugStructure          `json:"drugStructure"`
	Synonyms               []string                `json:"synonyms"`
	Pharmacology           *DrugPharmacology       `json:"pharmacology"`
	ExperimentalProperties *ExperimentalProperties `json:"experimentalProperties"`
}
