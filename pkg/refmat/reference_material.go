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
type ReferenceMap map[string][]*ReferenceMaterial

// BuildReferenceMap builds a map of all terms to their reference materials.
func BuildReferenceMap(refs []*ReferenceMaterial) ReferenceMap {
	m := make(map[string][]*ReferenceMaterial)
	for _, ref := range refs {
		for _, term := range ref.Terms {
			m[term] = append(m[term], ref)
		}
	}
	return m
}
