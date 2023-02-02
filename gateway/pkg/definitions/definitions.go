package definitions

type Gateway interface {
	PushEvent(Event) Void
}

type Event struct {
	ProjectIDs []string `json:"projectIDs"`
	Type       string   `json:"type"`
	Payload    string   `json:"payload"`
}

type Void struct{}
