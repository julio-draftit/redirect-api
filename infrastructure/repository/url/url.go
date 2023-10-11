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
		SELECT u.id, url, pixel 
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
		if err := row.Scan(&url.ID, &url.Url, &url.Pixel); err != nil {
			return nil, err
		}
	}

	if url.ID == 0 {
		return nil, nil
	}

	return &url, err
}
