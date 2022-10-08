package repository

import (
	"time"
)

// Item represents an item in our system
type Item struct {
	ID            *int64     `db:"id"`
	Quantity      int32      `db:"quantity"`
	Name          string     `db:"name"`
	ReservationID *int64     `db:"reservation_id"`
	ReservedAt    *time.Time `db:"reserved_at"`
	CreatedAt     *time.Time `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`
}
