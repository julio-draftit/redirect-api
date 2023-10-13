package user

import (
	"context"

	coreUser "github.com/Projects-Bots/redirect/internal/core/user"
	entity "github.com/Projects-Bots/redirect/internal/core/user"
)

type UserService struct {
	repository entity.NewUserRepositoryInterface
}

func NewUserService(repository entity.NewUserRepositoryInterface) *UserService {
	return &UserService{repository: repository}
}

func (s UserService) SelectUser(ctx context.Context, userID int) (*entity.User, error) {
	return s.repository.SelectUser(ctx, userID)
}
func (s *UserService) ListUser(ctx context.Context) ([]entity.User, error) {
	return s.repository.ListUser(ctx)
}
func (s *UserService) AddUser(ctx context.Context, user entity.User) (*entity.User, error) {
	return s.repository.AddUser(ctx, user)
}
func (s *UserService) Auth(ctx context.Context, user entity.User) (*entity.User, error) {
	return s.repository.Auth(ctx, user)
}
func (s *UserService) UpdateUser(ctx context.Context, id int, user entity.User) (*entity.User, error) {
	return s.repository.UpdateUser(ctx, id, user)
}
func (s *UserService) DeleteUser(ctx context.Context, id int) (*coreUser.User, error) {
	return s.repository.DeleteUser(ctx, id)
}
