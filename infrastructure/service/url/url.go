package url

import (
	"context"
	entity "github.com/Projects-Bots/redirect/internal/core/url"
)

type UrlService struct {
	repository entity.NewUrlRepositoryInterface
}

func NewUrlService(repository entity.NewUrlRepositoryInterface) *UrlService {
	return &UrlService{repository: repository}
}

func (s UrlService) GetUrl(ctx context.Context, url string) (*entity.Url, error) {
	return s.repository.GetUrl(ctx, url)
}
