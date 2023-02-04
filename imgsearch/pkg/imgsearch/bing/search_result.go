package bing_imgsearch

type searchResult struct {
	_            struct{} `type:"structure"`
	Type         string   `json:"_type"`
	ReadLink     string   `json:"readLink"`
	WebSearchUrl string   `json:"webSearchUrl"`
	QueryContext struct {
		OriginalQuery           string `json:"originalQuery"`
		AlterationDisplayQuery  string `json:"alterationDisplayQuery"`
		AlterationOverrideQuery string `json:"alterationOverrideQuery"`
		AlterationMethod        string `json:"alterationMethod"`
		AlterationType          string `json:"alterationType"`
	} `json:"queryContext"`
	TotalEstimatedMatches int64 `json:"totalEstimatedMatches"`
	NextOffset            int64 `json:"nextOffset"`
	CurrentOffset         int64 `json:"currentOffset"`
	Value                 []struct {
		WebSearchUrl       string `json:"webSearchUrl"`
		Name               string `json:"name"`
		ThumbnailUrl       string `json:"thumbnailUrl"`
		ContentUrl         string `json:"contentUrl"`
		HostPageUrl        string `json:"hostPageUrl"`
		EncodingFormat     string `json:"encodingFormat"`
		ContentSize        string `json:"contentSize"`
		HostPageDisplayUrl string `json:"hostPageDisplayUrl"`
		Width              int    `json:"width"`
		Height             int    `json:"height"`
		Thumbnail          struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"thumbnail"`
		ImageInsightsToken string `json:"imageInsightsToken"`
		InsightsMetadata   struct {
			PagesIncludingCount int `json:"pagesIncludingCount"`
			AvailableSizesCount int `json:"availableSizesCount"`
		} `json:"insightsMetadata"`
		ImageId     string `json:"imageId"`
		AccentColor string `json:"accentColor"`
	} `json:"value"`
	QueryExpansions []struct {
		Text         string `json:"text"`
		DisplayText  string `json:"displayText"`
		WebSearchUrl string `json:"webSearchUrl"`
		SearchLink   string `json:"searchLink"`
		Thumbnail    struct {
			ThumbnailUrl string `json:"thumbnailUrl"`
		} `json:"thumbnail"`
	} `json:"queryExpansions"`
	PivotSuggestions []struct {
		Pivot       string `json:"pivot"`
		Suggestions []struct {
			DisplayText  string `json:"displayText"`
			WebSearchUrl string `json:"webSearchUrl"`
			SearchLink   string `json:"searchLink"`
			Thumbnail    struct {
				ThumbnailUrl string `json:"thumbnailUrl"`
			} `json:"thumbnail"`
		} `json:"suggestions"`
	} `json:"pivotSuggestions"`
	RelatedSearches []struct {
		Text         string `json:"text"`
		DisplayText  string `json:"displayText"`
		WebSearchUrl string `json:"webSearchUrl"`
		SearchLink   string `json:"searchLink"`
		Thumbnail    struct {
			ThumbnailUrl string `json:"thumbnailUrl"`
		} `json:"thumbnail"`
	} `json:"relatedSearches"`
}
