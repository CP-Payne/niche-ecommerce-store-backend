// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

type Category struct {
	ID          uuid.UUID
	Name        string
	Description sql.NullString
}

type Product struct {
	ID             uuid.UUID
	Name           string
	Description    sql.NullString
	Price          string
	Brand          sql.NullString
	Sku            string
	StockQuantity  int32
	CategoryID     uuid.UUID
	ImageUrl       sql.NullString
	ThumbnailUrl   sql.NullString
	Specifications pqtype.NullRawMessage
	Variants       pqtype.NullRawMessage
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Review struct {
	ID         uuid.UUID
	Title      sql.NullString
	ReviewText sql.NullString
	Rating     int32
	ProductID  uuid.UUID
	UserID     uuid.UUID
	Deleted    bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type User struct {
	ID             uuid.UUID
	Name           sql.NullString
	Email          string
	HashedPassword string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
