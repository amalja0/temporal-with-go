package postgresrepository

import (
	"context"
	"database/sql"
	"sales-record-orchestration/internal/domain"
	"sales-record-orchestration/internal/ports"
	"strconv"
)

type postgresRepository struct {
	db *sql.DB
}

func InitPostgresRepository(db *sql.DB) ports.PostgresRepository {
	return &postgresRepository{
		db: db,
	}
}

func baseQuery() string {
	return `select
			s.id,
			s.ship_date,
			s.ship_mode,
			s.customer_name,
			s.quantity,
			s.sales_amount,
			s.discount,
			s.profit,
			s.profit_ratio,
			s.number_of_record,
			s.order_id,
			s.order_date,
			s.location_id,
			s.product_id,
			s.segment_id,
			concat(l.city, ', ', l.state, ', ', l.postal_code) as address,
			p.product_name ,
			sg.segment_name 
		from
			sales s
		left join locations l on
			l.id = s.location_id
		left join products p on
			p.id = s.product_id
		left join segments sg on
			sg.id = s.segment_id
		where
			1 = 1
	`
}

func (p *postgresRepository) FetchSales(ctx context.Context, queryParams *domain.SalesQueryParams) ([]domain.SaleWithDetail, error) {
	baseQuery := baseQuery()
	var args []any
	argIndex := 1

	if queryParams.OrderDate != nil && *queryParams.OrderDate != "" {
		baseQuery += ` and date(order_date) = $` + strconv.Itoa(argIndex)
		args = append(args, queryParams.OrderDate)
		argIndex++
	}

	rows, err := p.db.Query(baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sales []domain.SaleWithDetail
	for rows.Next() {
		var sale domain.SaleWithDetail

		if err := rows.Scan(
			&sale.ID,
			&sale.ShipDate,
			&sale.ShipMode,
			&sale.CustomerName,
			&sale.Qty,
			&sale.SaleAmount,
			&sale.Discount,
			&sale.Profit,
			&sale.ProfitRatio,
			&sale.NumberOfRecord,
			&sale.OrderID,
			&sale.OrderDate,
			&sale.LocationID,
			&sale.ProductID,
			&sale.SegmentID,
			&sale.Address,
			&sale.ProductName,
			&sale.SegmentName,
		); err != nil {
			return nil, err
		}
		sales = append(sales, sale)
	}
	
	return sales, nil
}

func (p *postgresRepository) FetchSale(ctx context.Context, queryParams domain.SaleQueryParams) (*domain.SaleWithDetail, error) {
	baseQuery := baseQuery()
	var args []any
	argIndex := 1

	if queryParams.SaleID != nil && *queryParams.SaleID != "" {
		baseQuery += ` or s.id = $` + strconv.Itoa(argIndex)
		args = append(args, queryParams.SaleID)
		argIndex++
	}

	row := p.db.QueryRow(baseQuery, args...)

	var sale domain.SaleWithDetail
	if err := row.Scan(
		&sale.ID,
		&sale.ShipDate,
		&sale.ShipMode,
		&sale.CustomerName,
		&sale.Qty,
		&sale.SaleAmount,
		&sale.Discount,
		&sale.Profit,
		&sale.ProfitRatio,
		&sale.NumberOfRecord,
		&sale.OrderID,
		&sale.OrderDate,
		&sale.LocationID,
		&sale.ProductID,
		&sale.SegmentID,
		&sale.Address,
		&sale.ProductName,
		&sale.SegmentName,
	); err != nil {
		return nil, err
	}

	return &sale, nil
}
