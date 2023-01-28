package refmat

type ReferenceMaterial struct {
	Terms  []string `json:"terms"`
	Images []string `json:"images"`
}

type ReferenceMap map[string]*ReferenceMaterial

func BuildReferenceMap(refs []*ReferenceMaterial) ReferenceMap {
	m := make(map[string]*ReferenceMaterial)
	for _, ref := range refs {
		for _, term := range ref.Terms {
			m[term] = ref
		}
	}
	return m
}

var TestReferenceMaterials = []*ReferenceMaterial{{
	Terms: []string{
		"vertebral arch",
	},
	Images: []string{
		"https://refmat.nyc3.digitaloceanspaces.com/lumbar-vertebra-vertebral-arch-superior-view-745x550.png",
		"https://refmat.nyc3.digitaloceanspaces.com/General-Structure-of-a-Vertebrae.jpg",
	},
}, {
	Terms: []string{
		"ligamentum flavum",
		"ligamentum",
		"flavum",
	},
	Images: []string{
		"https://refmat.nyc3.digitaloceanspaces.com/ligamentum-flavum-1024x670.jpg",
		"https://refmat.nyc3.digitaloceanspaces.com/LigamentumFlavum.png",
	},
}, {
	Terms: []string{
		"facet joint",
		"facet",
	},
	Images: []string{
		"https://refmat.nyc3.digitaloceanspaces.com/facet_joints_related_spine_structures_shutterstock_157672247.jpg",
		"https://refmat.nyc3.digitaloceanspaces.com/Thoracic-Facet-Syndrome.jpg",
	},
}}
