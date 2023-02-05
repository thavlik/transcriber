package imgsearch

type Result struct {
	Images          []*Image `json:"images"`
	QueryExpansions []string `json:"queryExpansions"`
}
