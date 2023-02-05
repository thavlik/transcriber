package definitions

type PharmaSeer interface {
	GetDrugDetails(GetDrugDetails) DrugDetails
}

type GetDrugDetails struct {
	Query string `json:"query"`
	Force bool   `json:"force"`
}

type DrugStructure struct {
	Image string `json:"image"`
	PDB   string `json:"pdb"`
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
	Indication           string   `json:"indication"`
	AssociatedConditions []string `json:"associatedConditions,omitempty"`
	Pharmacodynamics     string   `json:"pharmacodynamics"`
	MechanismOfAction    string   `json:"mechanismOfAction"`
	Absorption           string   `json:"absorption"`
	VolumeOfDistribution string   `json:"volumeOfDistribution"`
	ProteinBinding       string   `json:"proteinBinding"`
	Metabolism           string   `json:"metabolism"`
	RouteOfElimination   string   `json:"routeOfElimination"`
	HalfLife             string   `json:"halfLife"`
	Clearance            string   `json:"clearance"`
	Toxicity             string   `json:"toxicity"`
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

type DrugDetails struct {
	Summary                 string            `json:"summary"`
	BrandNames              []string          `json:"brandNames"`
	GenericName             string            `json:"genericName"`
	DrugBankAccessionNumber string            `json:"drugBankAccessionNumber"`
	Background              string            `json:"background"`
	Groups                  []string          `json:"groups"`
	Structure               *DrugStructure    `json:"drugStructure"`
	Weight                  *Weight           `json:"weight"`
	ChemicalFormula         string            `json:"chemicalFormula"`
	Synonyms                []string          `json:"synonyms"`
	ExternalIDs             []string          `json:"externalIDs"`
	Pharmacology            *DrugPharmacology `json:"pharmacology"`
	References              *References       `json:"references"`
}
