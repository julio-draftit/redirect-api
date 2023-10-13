package url

import (
	"context"
	"database/sql"

	entity "github.com/Projects-Bots/redirect/internal/core/url"
)

type UrlRepository struct {
	db *sql.DB
}

func NewUrlRepository(db *sql.DB) *UrlRepository {
	return &UrlRepository{
		db: db,
	}
}

func (r *UrlRepository) GetUrl(ctx context.Context, slug string) (*entity.Url, error) {
	sql := `
		SELECT u.id, url, pixel, random
		FROM urls u
		INNER JOIN users ON users.id = u.user_id 
		WHERE u.url = ?
		AND u.deleted_at IS NULL
		AND users.deleted_at IS NULL`

	row, err := r.db.Query(sql, slug)
	if err != nil {
		return nil, err
	}

	var url entity.Url
	if row.Next() {
		if err := row.Scan(&url.ID, &url.Url, &url.Pixel, &url.Random); err != nil {
			return nil, err
		}
	}

	if url.ID == 0 {
		return nil, nil
	}

	return &url, err
}
func (r *UrlRepository) GetAllUrlsByUser(ctx context.Context, userID int) ([]entity.Url, error) {
	sqlURLs := `
		SELECT u.id, u.user_id, u.name, u.url, u.pixel, u.random, u.created_at, u.updated_at
		FROM urls u
		WHERE u.user_id = ? AND u.deleted_at IS NULL
	`

	rows, err := r.db.QueryContext(ctx, sqlURLs, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []entity.Url
	for rows.Next() {
		var url entity.Url
		if err := rows.Scan(&url.ID, &url.UserID, &url.Name, &url.Url, &url.Pixel, &url.Random, &url.CreatedAt, &url.UpdatedAt); err != nil {
			return nil, err
		}

		const sqlRedirects = `
			SELECT r.url, r.hits, r.limit_hits
			FROM redirects r
			WHERE r.url_id = ? AND r.deleted_at IS NULL
		`
		redirectRows, err := r.db.QueryContext(ctx, sqlRedirects, url.ID)
		if err != nil {
			return nil, err
		}

		url.Redirects = []entity.Redirect{}
		for redirectRows.Next() {
			var redirect entity.Redirect
			if err := redirectRows.Scan(&redirect.URL, &redirect.Hits, &redirect.LimitHits); err != nil {
				redirectRows.Close()
				return nil, err
			}
			url.Redirects = append(url.Redirects, redirect)
		}
		redirectRows.Close()

		urls = append(urls, url)
	}

	return urls, nil
}

func (r *UrlRepository) AddUrl(ctx context.Context, url entity.Url) (*entity.Url, error) {
	sqlInsert := `
		INSERT INTO urls (user_id, name, url, pixel, random) values (?, ?, ?, ?, ?)
	`

	stmt, err := r.db.PrepareContext(ctx, sqlInsert)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(url.UserID, url.Name, url.Url, url.Pixel, url.Random)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	url.ID = int(id)

	for _, redirect := range url.Redirects {
		sqlInsertRedirect := `
			INSERT INTO redirects (url_id, url, hits, limit_hits) values (?, ?, ?, ?)
		`
		_, err = r.db.ExecContext(ctx, sqlInsertRedirect, url.ID, redirect.URL, redirect.Hits, redirect.LimitHits)
		if err != nil {
			return nil, err
		}
	}

	sqlFetch := `
		SELECT id, user_id, name, url, pixel, random, created_at, updated_at, deleted_at
		FROM urls WHERE id = ?
	`
	err = r.db.QueryRowContext(ctx, sqlFetch, url.ID).Scan(
		&url.ID, &url.UserID, &url.Name, &url.Url, &url.Pixel, &url.Random,
		&url.CreatedAt, &url.UpdatedAt, &url.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	return &url, nil
}

func (r *UrlRepository) UpdateUrl(ctx context.Context, id int, url entity.Url) (*entity.Url, error) {
	sqlUpdate := `
        UPDATE urls SET name = ?, url = ?, pixel = ?, random = ?, updated_at = NOW() WHERE id = ?
    `
	_, err := r.db.ExecContext(ctx, sqlUpdate, url.Name, url.Url, url.Pixel, url.Random, id)
	if err != nil {
		return nil, err
	}

	_, err = r.db.ExecContext(ctx, "DELETE FROM redirects WHERE url_id = ?", id)
	if err != nil {
		return nil, err
	}

	for _, redirect := range url.Redirects {
		_, err := r.db.ExecContext(ctx, "INSERT INTO redirects (url_id, url, hits, limit_hits) VALUES (?, ?, ?, ?)", id, redirect.URL, redirect.Hits, redirect.LimitHits)
		if err != nil {
			return nil, err
		}
	}

	sqlSelect := `
        SELECT id, user_id, name, url, pixel, random, created_at, updated_at, deleted_at FROM urls WHERE id = ?
    `
	row := r.db.QueryRowContext(ctx, sqlSelect, id)

	updatedUrl := entity.Url{}
	err = row.Scan(&updatedUrl.ID, &updatedUrl.UserID, &updatedUrl.Name, &updatedUrl.Url, &updatedUrl.Pixel, &updatedUrl.Random, &updatedUrl.CreatedAt, &updatedUrl.UpdatedAt, &updatedUrl.DeletedAt)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, "SELECT url, hits, limit_hits FROM redirects WHERE url_id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var redirect entity.Redirect
		if err := rows.Scan(&redirect.URL, &redirect.Hits, &redirect.LimitHits); err != nil {
			return nil, err
		}
		updatedUrl.Redirects = append(updatedUrl.Redirects, redirect)
	}

	return &updatedUrl, nil
}

func (r *UrlRepository) DeleteUrl(ctx context.Context, id int) (*entity.Url, error) {

	query := "UPDATE urls SET deleted_at = NOW() WHERE id = ?"

	url, err := r.GetUrlById(ctx, id)
	if err != nil {
		return nil, err
	}

	_, err = r.db.ExecContext(ctx, query, id)
	if err != nil {
		return nil, err
	}

	return url, nil
}
func (r *UrlRepository) GetUrlById(ctx context.Context, id int) (*entity.Url, error) {
	var url entity.Url
	query := "SELECT id, user_id, name, url, pixel, random, created_at, updated_at, deleted_at FROM urls WHERE id = ?"
	err := r.db.QueryRowContext(ctx, query, id).Scan(&url.ID, &url.UserID, &url.Name, &url.Url, &url.Pixel, &url.Random, &url.CreatedAt, &url.UpdatedAt, &url.DeletedAt)
	if err != nil {
		return nil, err
	}

	redirectQuery := "SELECT url, hits, limit_hits FROM redirects WHERE url_id = ? AND deleted_at IS NULL"
	rows, err := r.db.QueryContext(ctx, redirectQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var redirects []entity.Redirect
	for rows.Next() {
		var redirect entity.Redirect
		if err := rows.Scan(&redirect.URL, &redirect.Hits, &redirect.LimitHits); err != nil {
			return nil, err
		}
		redirects = append(redirects, redirect)
	}

	url.Redirects = redirects

	return &url, nil
}
