package comprehend

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

// Entity represents a detected key term or phrase.
// The Type field is dependent on the Comprehend service used.
type Entity struct {
	Text  string  `json:"text"`  // text of entity
	Type  string  `json:"type"`  // type of entity
	Score float64 `json:"score"` // confidence score, 0.0-1.0
}

// Hash returns the unique md5 hash for the entity.
// This is used to identify duplicate entities and
// does not include the Score field.
func (e *Entity) Hash() string {
	h := md5.New()
	h.Write([]byte(e.Text))
	h.Write([]byte(e.Type))
	raw := h.Sum(nil)
	return hex.EncodeToString(raw[:])
}

func (e *Entity) String() string {
	return fmt.Sprintf("%s (%s)", e.Text, e.Type)
}
