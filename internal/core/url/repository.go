package url

import "context"

type NewUrlRepositoryInterface interface {
	GetUrl(ctx context.Context, url string) (*Url, error)
	GetAllUrlsByUser(ctx context.Context, userID int) ([]Url, error)
	AddUrl(ctx context.Context, url Url) (*Url, error)
	UpdateUrl(ctx context.Context, id int, url Url) (*Url, error)
	DeleteUrl(ctx context.Context, id int) (*Url, error)
}
