package user

import (
	"context"
	"database/sql"

	coreUser "github.com/Projects-Bots/redirect/internal/core/user"
	entity "github.com/Projects-Bots/redirect/internal/core/user"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) SelectUser(ctx context.Context, user_id int) (*entity.User, error) {
	sql := `
		SELECT id, name, email, admin, deleted_at 
		FROM users where id = ?`

	row, err := r.db.Query(sql, user_id)
	if err != nil {
		return nil, err
	}

	var user entity.User
	if row.Next() {
		if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Admin, &user.DeletedAt); err != nil {
			return nil, err
		}
	}

	if user.ID == 0 {
		return nil, nil
	}

	return &user, err
}
func (r *UserRepository) Auth(ctx context.Context, user entity.User) (*entity.User, error) {
	sqlQuery := `
		SELECT id, name, email, admin, created_at, updated_at, deleted_at 
		FROM users 
		WHERE email = ? AND password = ? AND deleted_at IS NULL`

	var fetchedUser entity.User
	r.db.QueryRowContext(ctx, sqlQuery, user.Email, user.Password).Scan(
		&fetchedUser.ID,
		&fetchedUser.Name,
		&fetchedUser.Email,
		&fetchedUser.Admin,
		&fetchedUser.CreatedAt,
		&fetchedUser.UpdatedAt,
		&fetchedUser.DeletedAt,
	)

	if fetchedUser.ID == 0 {
		return nil, nil
	}

	return &fetchedUser, nil
}
func (r *UserRepository) ListUser(ctx context.Context) ([]entity.User, error) {
	sql := `
		SELECT id, name, email, admin, created_at, updated_at, deleted_at 
		FROM users
	`

	rows, err := r.db.QueryContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Admin, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) AddUser(ctx context.Context, user entity.User) (*entity.User, error) {
	sqlInsert := `
		INSERT INTO users (name, email, password, admin) values (?, ?, ?, ?)
	`

	stmt, err := r.db.PrepareContext(ctx, sqlInsert)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(user.Name, user.Email, user.Password, user.Admin)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	user.ID = int(id)

	sqlFetch := `
		SELECT id, name, email, admin, deleted_at 
		FROM users where id = ?
	`
	err = r.db.QueryRowContext(ctx, sqlFetch, user.ID).Scan(
		&user.ID, &user.Name, &user.Email, &user.Admin, &user.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, id int, user entity.User) (*entity.User, error) {
	sqlUpdate := `
        UPDATE users SET name = ?, email = ?, admin = ?, updated_at = NOW() WHERE id = ?
    `
	_, err := r.db.ExecContext(ctx, sqlUpdate, user.Name, user.Email, user.Admin, id)
	if err != nil {
		return nil, err
	}

	sqlSelect := `
        SELECT id, name, email, admin, created_at, updated_at, deleted_at FROM users WHERE id = ?
    `
	row := r.db.QueryRowContext(ctx, sqlSelect, id)

	updatedUser := coreUser.User{}
	err = row.Scan(&updatedUser.ID, &updatedUser.Name, &updatedUser.Email, &updatedUser.Admin, &updatedUser.CreatedAt, &updatedUser.UpdatedAt, &updatedUser.DeletedAt)
	if err != nil {
		return nil, err
	}
	return &updatedUser, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int) (*coreUser.User, error) {

	query := "UPDATE users SET deleted_at = NOW() WHERE id = ?"

	user, err := r.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	_, err = r.db.ExecContext(ctx, query, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
func (r *UserRepository) GetUserById(ctx context.Context, id int) (*coreUser.User, error) {
	var user coreUser.User
	query := "SELECT id, name, email, admin, created_at, updated_at, deleted_at FROM users WHERE id = ?"
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Admin, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
