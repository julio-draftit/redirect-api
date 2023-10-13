package url

type Url struct {
	ID        int        `json:"id"`
	UserID    int        `json:"user_id"`
	Name      string     `json:"name"`
	Url       string     `json:"url"`
	Pixel     string     `json:"pixel"`
	Random    *bool      `json:"random"`
	Redirects []Redirect `json:"redirects"`
	CreatedAt string     `json:"created_at"`
	UpdatedAt string     `json:"updated_at"`
	DeletedAt *string    `json:"deleted_at"`
}

type Redirect struct {
	URL       string `json:"url"`
	Hits      int    `json:"hits"`
	LimitHits int    `json:"limit_hits"`
}
