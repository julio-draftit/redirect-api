package access

import "context"

type NewAcesssRepositoryInterface interface {
	Save(ctx context.Context, redirectID int) (*Access, error)
}
