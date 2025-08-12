package temporal

import (
	"context"
	"sales-record-orchestration/internal/domain"
	"sales-record-orchestration/internal/ports"
)

type temporalActivity struct {
	postgresRepo   ports.PostgresRepository
	kafkaProducer  ports.KafkaProducer
	clickhouseRepo ports.ClickHouseRepository
}

func InitActivities(
	postgresRepo ports.PostgresRepository,
	kafkaProducer ports.KafkaProducer,
	clickhouseRepo ports.ClickHouseRepository,
) ports.TemporalActivity {
	return &temporalActivity{
		postgresRepo:   postgresRepo,
		kafkaProducer:  kafkaProducer,
		clickhouseRepo: clickhouseRepo,
	}
}

func (t *temporalActivity) FetchSalesActivity(ctx context.Context, queryParams *domain.SalesQueryParams) ([]domain.SaleWithDetail, error) {
	return t.postgresRepo.FetchSales(ctx, queryParams)
}

func (t *temporalActivity) PublishSalesActivity(ctx context.Context, sales []domain.SaleWithDetail) error {
	return t.kafkaProducer.PublishSales(ctx, "sales", sales)
}

func (t *temporalActivity) ProcessSalesActivity(ctx context.Context, records []domain.OrderRecord) error {
	return t.clickhouseRepo.StoreSales(ctx, records)
}
