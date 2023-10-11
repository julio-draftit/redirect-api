package access

import (
	"context"
	entity "github.com/Projects-Bots/redirect/internal/core/access"
)

type AccessService struct {
	repository entity.NewAcesssRepositoryInterface
}

func NewAccessService(repository entity.NewAcesssRepositoryInterface) *AccessService {
	return &AccessService{repository: repository}
}

func (r AccessService) Save(ctx context.Context, redirectID int) (*entity.Access, error) {
	return r.repository.Save(ctx, redirectID)
}
