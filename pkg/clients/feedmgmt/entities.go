package feedmgmt

type Feed struct {
	URL      string `json:"url"`
	Provider string `json:"provider"`
	Category string `json:"category"`
}

type Feeds []Feed
