package domain

import (
	"database/sql"
	"time"
)

type Log struct {
	CreatedBy string         `db:"created_by"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedBy string         `db:"updated_by"`
	UpdatedAt time.Time      `db:"updated_at"`
	DeletedBy sql.NullString `db:"deleted_by"`
	DeletedAt time.Time      `db:"deleted_at"`
}
