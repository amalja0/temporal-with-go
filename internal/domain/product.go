package domain

import "github.com/google/uuid"

type Product struct {
	ID            uuid.UUID `db:"id" json:"id"`
	ProductName   string    `db:"product_name" json:"product_name"`
	Manufacturer  string    `db:"manufacturer" json:"manufacturer"`
	BasePrice     float32   `db:"base_price" json:"base_price"`
	CategoryID    uuid.UUID `db:"category_id" json:"category_id"`
	SubCategoryID uuid.UUID `db:"sub_category_id" json:"sub_category_id"`
	*Log
}
