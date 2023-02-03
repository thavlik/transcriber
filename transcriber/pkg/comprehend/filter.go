package comprehend

import "github.com/thavlik/transcriber/base/pkg/base"

type Filter struct {
	IncludeTypes []string `json:"includeTypes"` // if not empty, only entities of these types are returned
	ExcludeTypes []string `json:"excludeTypes"` // if not empty, entities of these types are filtered out
	IncludeTerms []string `json:"includeTerms"` // if not empty, only entities matching these case-insensitive terms are returned
	ExcludeTerms []string `json:"excludeTerms"` // if not empty, entities matching these case-insensitive terms are filtered out
}

func (f *Filter) Matches(entity *Entity) bool {
	if f == nil {
		return true
	}
	omit := (len(f.IncludeTypes) > 0 && !base.Contains(f.IncludeTypes, entity.Type)) ||
		(len(f.ExcludeTypes) > 0 && base.Contains(f.ExcludeTypes, entity.Type)) ||
		(len(f.IncludeTerms) > 0 && !base.Contains(f.IncludeTerms, entity.Text)) ||
		(len(f.ExcludeTerms) > 0 && base.Contains(f.ExcludeTerms, entity.Text))
	return !omit
}
