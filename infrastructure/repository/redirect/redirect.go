package redirect

import (
	"context"
	"database/sql"

	entity "github.com/Projects-Bots/redirect/internal/core/redirect"
)

type RedirectRepository struct {
	db *sql.DB
}

func NewRedirectRepository(db *sql.DB) *RedirectRepository {
	return &RedirectRepository{
		db: db,
	}
}

func (r RedirectRepository) GetUrl(ctx context.Context, urlID int, random bool) (*entity.Redirect, error) {
	baseSql := `
		SELECT r.id, r.url_id, r.url, r.hits, r.limit_hits
		FROM redirects r 
		INNER JOIN urls u ON u.id = r.url_id 
		INNER JOIN users us ON us.id = u.user_id 
		WHERE r.url_id = ?
		AND hits < limit_hits 
		AND r.deleted_at IS NULL
		AND u.deleted_at IS NULL
		AND us.deleted_at IS NULL
	`

	orderSql := ""
	if random {
		orderSql = "ORDER BY RAND()"
	}

	limitSql := "LIMIT 1"

	finalSql := baseSql + orderSql + limitSql

	row, err := r.db.Query(finalSql, urlID)
	if err != nil {
		return nil, err
	}

	var redirect entity.Redirect
	if row.Next() {
		if err := row.Scan(&redirect.ID, &redirect.UrlID, &redirect.Url, &redirect.Hits, &redirect.LimitHits); err != nil {
			return nil, err
		}
	}

	if redirect.ID == 0 {
		return nil, nil
	}

	return &redirect, err
}

func (r RedirectRepository) UpdateUrl(ctx context.Context, id int) error {
	sql := `
		UPDATE redirects SET hits = (hits +1)
		WHERE id = ?	
	`

	row, err := r.db.Prepare(sql)
	if err != nil {
		return err
	}

	defer row.Close()

	if _, err = row.Exec(id); err != nil {
		return err
	}

	return nil
}

func (r RedirectRepository) ResetHitsUrl(ctx context.Context, urlID int) error {
	sql := `
		UPDATE redirects SET hits = 0
		WHERE url_id = ?
	`

	row, err := r.db.Prepare(sql)
	if err != nil {
		return err
	}

	defer row.Close()

	if _, err = row.Exec(urlID); err != nil {
		return err
	}

	return nil
}
