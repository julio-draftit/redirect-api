package access

import (
	"context"
	"database/sql"
	entity "github.com/Projects-Bots/redirect/internal/core/access"
)

type AccessRepository struct {
	db *sql.DB
}

func NewAccessRepository(db *sql.DB) *AccessRepository {
	return &AccessRepository{
		db: db,
	}
}

func (r AccessRepository) Save(ctx context.Context, redirectID int) (*entity.Access, error) {
	sql := `
		INSERT INTO accesses (redirect_id) values(?)	
	`

	row, err := r.db.Prepare(sql)
	if err != nil {
		return nil, err
	}

	defer row.Close()

	result, err := row.Exec(redirectID)
	if err != nil {
		return nil, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &entity.Access{
		ID:         int(lastID),
		RedirectID: redirectID,
	}, nil
}
