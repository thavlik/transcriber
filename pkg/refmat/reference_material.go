package refmat

// ReferenceMaterial is a reference material that can be used to explain a term
// that was found in a transcript.
type ReferenceMaterial struct {
	// Terms is a list of terms that can be used to find this reference material.
	// For example, if the reference material is about the ligamentum flavum,
	// then the terms could be "ligamentum flavum", "ligamentum", and "flavum".
	// The maximum number of words in a term is 2. Terms are not case sensitive.
	Terms []string `json:"terms"`

	// List of image URLs that can be used to display the reference material.
	// The first image in the list is the primary image.
	Images []string `json:"images"`
}

// ReferenceMap is a map of all terms to their reference materials.
type ReferenceMap map[string]*ReferenceMaterial

// BuildReferenceMap builds a map of all terms to their reference materials.
func BuildReferenceMap(refs []*ReferenceMaterial) ReferenceMap {
	m := make(map[string]*ReferenceMaterial)
	for _, ref := range refs {
		for _, term := range ref.Terms {
			m[term] = ref
		}
	}
	return m
}

// TestReferenceMaterials is a list of reference materials that can be used for
// testing. This list is not exhaustive. It is meant to be used for testing
// purposes only. This is a good test video: https://youtu.be/gd4-FV_lwSE
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
	},
	Images: []string{
		"https://refmat.nyc3.digitaloceanspaces.com/facet_joints_related_spine_structures_shutterstock_157672247.jpg",
		"https://refmat.nyc3.digitaloceanspaces.com/Thoracic-Facet-Syndrome.jpg",
	},
}}
