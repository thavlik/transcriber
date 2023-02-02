package comprehend

import "github.com/thavlik/transcriber/base/pkg/base"

// filter returns true if the type should be filtered.
// If includeTypes is not empty, the type must be in the list.
// If excludeTypes is not empty, the type must not be in the list.
func filter(
	ty string,
	includeTypes []string,
	excludeTypes []string,
) bool {
	return (len(includeTypes) > 0 && !base.Contains(includeTypes, ty)) || (len(excludeTypes) > 0 && base.Contains(excludeTypes, ty))
}
