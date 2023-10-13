package report

import "context"

type NewReportRepositoryInterface interface {
	Reports(ctx context.Context, report Report) (*Report, error)
}
