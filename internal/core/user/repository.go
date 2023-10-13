package user

import "context"

type NewUserRepositoryInterface interface {
	//Auth(ctx context.Context, url string) (*User, error)

	SelectUser(ctx context.Context, userID int) (*User, error)
	ListUser(ctx context.Context) ([]User, error)
	AddUser(ctx context.Context, url User) (*User, error)
	Auth(ctx context.Context, url User) (*User, error)
	UpdateUser(ctx context.Context, id int, url User) (*User, error)
	DeleteUser(ctx context.Context, id int) (*User, error)
}
