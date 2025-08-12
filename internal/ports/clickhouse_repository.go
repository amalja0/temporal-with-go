package ports

import (
	"context"
	"sales-record-orchestration/internal/domain"
)

type ClickHouseRepository interface {
	StoreSale(ctx context.Context, record domain.OrderRecord) error
	StoreSales(ctx context.Context, record []domain.OrderRecord) error
}
