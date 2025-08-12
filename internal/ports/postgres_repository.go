package ports

import (
	"context"
	"sales-record-orchestration/internal/domain"
)

type PostgresRepository interface {
	FetchSales(ctx context.Context, queryParams *domain.SalesQueryParams) ([]domain.SaleWithDetail, error)
	FetchSale(ctx context.Context, queryParams domain.SaleQueryParams) (*domain.SaleWithDetail, error)
}
