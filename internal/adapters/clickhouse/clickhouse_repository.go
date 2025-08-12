package clickhouserepository

import (
	"context"
	"fmt"
	"log"
	"sales-record-orchestration/internal/domain"
	"sales-record-orchestration/internal/ports"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type clickhouseRepository struct {
	db driver.Conn
}

func InitClickhouseRepository(db *driver.Conn) ports.ClickHouseRepository {
	return &clickhouseRepository{
		db: *db,
	}
}

func (c *clickhouseRepository) StoreSale(ctx context.Context, record domain.OrderRecord) error {
	query := `
		INSERT INTO realtime_order (
			id,
			sale_id,
			quantity,
			sale_amount,
			discount,
			profit,
			profit_ratio,
			order_id,
			order_date,
			location_id,
			product_id,
			segment_id,
			product_name,
			segment_name,
			created_at
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)
	`

	err := c.db.Exec(
		context.Background(),
		query,
		record.Id,
		record.SaleId,
		record.Quantity,
		record.SaleAmount,
		record.Discount,
		record.Profit,
		record.ProfitRatio,
		record.OrderId,
		record.OrderDate,
		record.LocationId,
		record.ProductId,
		record.SegmentId,
		record.ProductName,
		record.SegmentName,
		record.CreatedAt,
	)
	if err != nil {
		fmt.Println("Fail to insert record:", err)
	}
	return err
}

func (c *clickhouseRepository) StoreSales(ctx context.Context, records []domain.OrderRecord) error {
	batch, err := c.db.PrepareBatch(ctx, "INSERT INTO realtime_order")
	if err != nil {
		log.Fatal("Failed to prepare batch:", err)
	}

	for _, record := range records {
		if err := batch.Append(
			record.Id,
			record.SaleId,
			record.Quantity,
			record.SaleAmount,
			record.Discount,
			record.Profit,
			record.ProfitRatio,
			record.OrderId,
			record.OrderDate,
			record.LocationId,
			record.ProductId,
			record.SegmentId,
			record.ProductName,
			record.SegmentName,
			record.CreatedAt,
		); err != nil {
			log.Fatal("Failed to append row:", err)
		}
	}

	if err := batch.Send(); err != nil {
		log.Fatal("Failed to send batch:", err)
	}

	return err
}
