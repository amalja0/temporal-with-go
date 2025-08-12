package ports

import (
	"context"
	"sales-record-orchestration/internal/domain"
)

type TemporalActivity interface {
	FetchSalesActivity(ctx context.Context, queryParams *domain.SalesQueryParams) ([]domain.SaleWithDetail, error)
	PublishSalesActivity(ctx context.Context, sales []domain.SaleWithDetail) error
	ProcessSalesActivity(ctx context.Context, records []domain.OrderRecord) error
}
