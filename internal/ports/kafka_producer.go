package ports

import (
	"context"
	"sales-record-orchestration/internal/domain"
)

type KafkaProducer interface {
	PublishSale(ctx context.Context, topic string, sale domain.SaleWithDetail) error
	PublishSales(ctx context.Context, topic string, sales []domain.SaleWithDetail) error
}
