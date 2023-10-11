package url

import "context"

type NewUrlRepositoryInterface interface {
	GetUrl(ctx context.Context, url string) (*Url, error)
}
