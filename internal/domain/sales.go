package domain

import (
	"time"

	"github.com/google/uuid"
)

type Sale struct {
	ID             uuid.UUID `db:"id" json:"id"`
	ShipDate       time.Time `db:"ship_date" json:"ship_date"`
	ShipMode       string    `db:"ship_mode" json:"ship_mode"`
	CustomerName   string    `db:"customer_name" json:"customer_name"`
	Qty            int32     `db:"quantity" json:"quantity"`
	SaleAmount     float32   `db:"sale_amount" json:"sale_amount"`
	Discount       float32   `db:"discount" json:"discount"`
	Profit         float32   `db:"profit" json:"profit"`
	ProfitRatio    float32   `db:"profit_ratio" json:"profit_ratio"`
	NumberOfRecord int32     `db:"number_of_record" json:"number_of_record"`
	OrderID        string    `db:"order_id" json:"order_id"`
	OrderDate      time.Time `db:"order_date" json:"order_date"`
	LocationID     uuid.UUID `db:"location_id" json:"location_id"`
	ProductID      uuid.UUID `db:"product_id" json:"product_id"`
	SegmentID      uuid.UUID `db:"segment_id" json:"segment_id"`
	*Log
}

type SaleWithDetail struct {
	Sale
	Address     string `json:"address"`
	ProductName string `json:"product_name"`
	SegmentName string `json:"segment_name"`
}

type SaleQueryParams struct {
	SaleID *string `json:"sale_id"`
}

type SalesQueryParams struct {
	OrderDate *string `json:"order_date"`
}
