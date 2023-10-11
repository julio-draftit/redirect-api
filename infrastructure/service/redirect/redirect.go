package redirect

import (
	"context"
	entity "github.com/Projects-Bots/redirect/internal/core/redirect"
)

type RedirectService struct {
	repository entity.NewRedirectRepositoryInterface
}

func NewRedirectService(repository entity.NewRedirectRepositoryInterface) *RedirectService {
	return &RedirectService{repository: repository}
}

func (s RedirectService) GetUrl(ctx context.Context, urlID int) (*entity.Redirect, error) {
	return s.GetNextUrl(ctx, urlID)
}

func (s RedirectService) UpdateUrl(ctx context.Context, id int) error {
	return s.repository.UpdateUrl(ctx, id)
}

func (s RedirectService) ResetHitsUrl(ctx context.Context, urlID int) error {
	return s.repository.ResetHitsUrl(ctx, urlID)
}

func (s RedirectService) GetNextUrl(ctx context.Context, urlID int) (*entity.Redirect, error) {
	redirect, err := s.repository.GetUrl(ctx, urlID)
	if err != nil {
		return nil, err
	}

	if redirect == nil {
		err := s.repository.ResetHitsUrl(ctx, urlID)
		if err != nil {
			return nil, err
		}

		redirect, err = s.GetUrl(ctx, urlID)
		if err != nil {
			return nil, err
		}
	}

	return redirect, nil
}
