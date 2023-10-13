package report

import (
	"context"

	report "github.com/Projects-Bots/redirect/internal/core/report"
)

type ReportService struct {
	repository report.NewReportRepositoryInterface
}

func NewReportService(repository report.NewReportRepositoryInterface) *ReportService {
	return &ReportService{repository: repository}
}

func (s *ReportService) Reports(ctx context.Context, report report.Report) (*report.Report, error) {
	return s.repository.Reports(ctx, report)
}
