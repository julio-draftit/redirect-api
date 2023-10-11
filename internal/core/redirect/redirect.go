package redirect

type Redirect struct {
	ID        int    `json:"id"`
	UrlID     int    `json:"url_id"`
	Url       string `json:"url"`
	Hits      string `json:"hits"`
	LimitHits string `json:"limit_hits"`
}
