package report

import (
	"context"
	"database/sql"

	report "github.com/Projects-Bots/redirect/internal/core/report"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{
		db: db,
	}
}

func (r *ReportRepository) Reports(ctx context.Context, reports report.Report) (*report.Report, error) {
	sql := `
	SELECT ? as user_id,
    (SELECT COUNT(*) FROM urls WHERE user_id = ?) AS instancias,
    (SELECT COUNT(*) FROM redirects r JOIN urls u ON u.id = r.url_id WHERE u.user_id = ?) AS links,
    (SELECT CASE WHEN SUM(limit_hits) IS NULL THEN 0 ELSE SUM(limit_hits) end FROM redirects r JOIN urls u ON u.id = r.url_id WHERE u.user_id = ?) AS hits,
    (SELECT COUNT(*) FROM redirects r JOIN urls u ON u.id = r.url_id JOIN accesses a ON r.id = a.redirect_id WHERE u.user_id = ?) AS clicks;`

	row := r.db.QueryRow(sql, reports.UserID, reports.UserID, reports.UserID, reports.UserID, reports.UserID)

	reportResult := report.Report{}
	err := row.Scan(&reportResult.UserID, &reportResult.Instancias, &reportResult.Links, &reportResult.Hits, &reportResult.Clicks)
	if err != nil {
		return nil, err
	}

	return &reportResult, nil
}
