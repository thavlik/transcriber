package websearch

type Result struct {
	Name             string `json:"name"`
	URL              string `json:"url"`
	IsFamilyFriendly bool   `json:"isFamilyFriendly"`
}
