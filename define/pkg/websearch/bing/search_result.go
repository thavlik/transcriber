package bing_websearch

type webSearchResult struct {
	Type         string `json:"_type"`
	QueryContext struct {
		OriginalQuery string `json:"originalQuery"`
	} `json:"queryContext"`
	WebPages struct {
		WebSearchUrl          string `json:"webSearchUrl"`
		TotalEstimatedMatches int64  `json:"totalEstimatedMatches"`
		Value                 []struct {
			ID               string `json:"id"`
			Name             string `json:"name"`
			URL              string `json:"url"`
			IsFamilyFriendly bool   `json:"isFamilyFriendly"`
			DisplayURL       string `json:"displayUrl"`
			Snippet          string `json:"snippet"`
			DateLastCrawled  string `json:"dateLastCrawled"`
			Language         string `json:"language"`
			IsNavigational   bool   `json:"isNavigational"`
		} `json:"value"`
	} `json:"webPages"`
	RankingResponse struct {
		Mainline struct {
			Items []struct {
				ResultIndex int    `json:"resultIndex"`
				AnswerType  string `json:"answerType"`
				Value       struct {
					ID   string `json:"id"`
					Type string `json:"type"`
				} `json:"value"`
			} `json:"items"`
		} `json:"mainline"`
	} `json:"rankingResponse"`
}
