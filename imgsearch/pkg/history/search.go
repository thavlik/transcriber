package history

import (
	"encoding/json"
	"time"
)

type Search struct {
	ID        string
	Query     string
	UserID    string
	Timestamp time.Time
}

func (s *Search) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.AsMap())
}

func (s *Search) UnmarshalJSON(data []byte) error {
	var raw struct {
		ID        string `json:"id"`
		Query     string `json:"query"`
		UserID    string `json:"userID"`
		Timestamp int64  `json:"timestamp"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	s.ID = raw.ID
	s.Query = raw.Query
	s.UserID = raw.UserID
	s.Timestamp = time.Unix(0, raw.Timestamp)
	return nil
}

func (s *Search) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"id":        s.ID,
		"query":     s.Query,
		"userID":    s.UserID,
		"timestamp": s.Timestamp.UnixNano(),
	}
}

func (s *Search) AsMongoMap() map[string]interface{} {
	return map[string]interface{}{
		"_id":       s.ID,
		"query":     s.Query,
		"userID":    s.UserID,
		"timestamp": s.Timestamp.UnixNano(),
	}
}
