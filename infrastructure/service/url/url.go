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
func (s *UrlService) GetAllUrlsByUser(ctx context.Context, userID int) ([]entity.Url, error) {
	return s.repository.GetAllUrlsByUser(ctx, userID)
}
func (s *UrlService) AddUrl(ctx context.Context, url entity.Url) (*entity.Url, error) {
	return s.repository.AddUrl(ctx, url)
}
func (s *UrlService) UpdateUrl(ctx context.Context, id int, url entity.Url) (*entity.Url, error) {
	return s.repository.UpdateUrl(ctx, id, url)
}
func (s *UrlService) DeleteUrl(ctx context.Context, id int) (*entity.Url, error) {
	return s.repository.DeleteUrl(ctx, id)
}
