package comprehend

import "github.com/thavlik/transcriber/base/pkg/base"

type Filter struct {
	IncludeTypes []string `json:"includeTypes,omitempty"` // if not empty, only entities of these types are returned
	ExcludeTypes []string `json:"excludeTypes,omitempty"` // if not empty, entities of these types are filtered out
	IncludeTerms []string `json:"includeTerms,omitempty"` // if not empty, only entities matching these case-insensitive terms are returned
	ExcludeTerms []string `json:"excludeTerms,omitempty"` // if not empty, entities matching these case-insensitive terms are filtered out
}

func (f *Filter) Matches(entity *Entity) bool {
	if f == nil {
		return true
	}
	if (len(f.IncludeTypes) > 0 && !base.Contains(f.IncludeTypes, entity.Type)) ||
		(len(f.IncludeTerms) > 0 && !base.Contains(f.IncludeTerms, entity.Text)) {
		// we only want to include entities that match the include filters
		return false
	}
	if (len(f.ExcludeTypes) > 0 && base.Contains(f.ExcludeTypes, entity.Type)) ||
		(len(f.ExcludeTerms) > 0 && base.Contains(f.ExcludeTerms, entity.Text)) {
		// Always exclude these matched entities.
		// This allows an entity included based on its
		// Type to then be excluded based on its Text.
		return false
	}
	return true
}
