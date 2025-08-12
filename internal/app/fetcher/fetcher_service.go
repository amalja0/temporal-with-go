package fetcher

import (
	"context"
	"sales-record-orchestration/internal/domain"
	"sales-record-orchestration/internal/ports"
)

type FetcherService interface {
	FetchAndPublishSale(ctx context.Context, queryParams domain.SaleQueryParams) error
	FetchAndPublishSales(ctx context.Context, queryParams domain.SalesQueryParams) error
}

type fetcherService struct {
	pgRepository  ports.PostgresRepository
	kafkaProducer ports.KafkaProducer
}

func InitFetcherService(pgRepo ports.PostgresRepository, kafkaProducer ports.KafkaProducer) FetcherService {
	return &fetcherService{
		pgRepository:  pgRepo,
		kafkaProducer: kafkaProducer,
	}
}

func (f *fetcherService) FetchAndPublishSale(ctx context.Context, queryParams domain.SaleQueryParams) error {
	//TODO implement me
	panic("implement me")
}

func (f *fetcherService) FetchAndPublishSales(ctx context.Context, queryParams domain.SalesQueryParams) error {
	//TODO implement me
	panic("implement me")
}
