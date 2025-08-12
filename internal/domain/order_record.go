package domain

import (
	"time"

	"github.com/google/uuid"
)

type OrderRecord struct {
	Id          uuid.UUID `ch:"id"`
	SaleId      uuid.UUID `ch:"sale_id"`
	Quantity    int32     `ch:"quantity"`
	SaleAmount  float32   `ch:"sale_amount"`
	Discount    float32   `ch:"discount"`
	Profit      float32   `ch:"profit"`
	ProfitRatio float32   `ch:"profit_ratio"`
	OrderId     string    `ch:"order_id"`
	OrderDate   time.Time `ch:"order_date"`
	LocationId  uuid.UUID `ch:"location_id"`
	ProductId   uuid.UUID `ch:"product_id"`
	SegmentId   uuid.UUID `ch:"segment_id"`
	ProductName string    `ch:"product_name"`
	SegmentName string    `ch:"segment_name"`
	CreatedAt   time.Time `ch:"created_at"`
}
