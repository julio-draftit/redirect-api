package redirect

import "context"

type NewRedirectRepositoryInterface interface {
	GetUrl(ctx context.Context, urlID int) (*Redirect, error)
	UpdateUrl(ctx context.Context, id int) error
	ResetHitsUrl(ctx context.Context, urlID int) error
}
