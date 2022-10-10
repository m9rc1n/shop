package repository

import (
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/m9rc1n/shop/pkg/log"
)

const (
	ListQuery    = "SELECT id, name, quantity, created_at, updated_at, reserved_at, reservation_id FROM shop.items"
	GetByIdQuery = "SELECT id, name, quantity, created_at, updated_at, reserved_at, reservation_id FROM shop.items WHERE id = $1"
	ReserveQuery = "UPDATE shop.items SET (reservation_id, reserved_at) = ($1, NOW()) WHERE id = $2"
	CreateQuery  = "INSERT INTO shop.items (name, quantity) VALUES ($1, $2) RETURNING id"
)

// Repository define item repository methods interface
type Repository interface {
	// List returns the list of items
	List() ([]*Item, error)
	// Get an item by id
	Get(id int64) (*Item, error)
	// Create new item in the storage
	Create(name string, quantity int32) (int64, error)
	// Reserve an item in the storage
	Reserve(id, reservationId int64) error
}

// repository persists item in database
type repository struct {
	db     *sqlx.DB
	logger log.Logger
}

// NewRepository creates a new item repository
func NewRepository(db *sqlx.DB, log log.Logger) Repository {
	return repository{db: db, logger: log}
}

// List return slice of items base on request parameters
func (r repository) List() ([]*Item, error) {
	items := make([]*Item, 0)
	err := r.db.Select(
		&items,
		ListQuery,
	)
	if err != nil {
		r.logger.Errorf("Failed to get list of items %s", err)
		return nil, err
	}
	return items, nil
}

// Create item from name and quantity
func (r repository) Create(name string, quantity int32) (int64, error) {
	var id int64
	err := r.db.QueryRow(CreateQuery, name, quantity).Scan(&id)
	if err != nil {
		r.logger.Errorf("Failed to create item of name :%s, because of %s", name, err)
		return -1, err
	}
	return id, nil
}

// Reserve item from name and quantity
func (r repository) Reserve(id, reservationId int64) error {
	_, err := r.db.Exec(ReserveQuery, reservationId, id)
	if err != nil {
		r.logger.Errorf("Failed to reserve item of id :%d, because of %s", id, err)
		return err
	}
	return nil
}

// Get item by id
func (r repository) Get(id int64) (*Item, error) {
	item := new(Item)
	err := r.db.Get(item, GetByIdQuery, id)
	if err != nil {
		r.logger.Errorf("Failed to retrieve item with id :%d, because of %s", id, err)
		return nil, err
	}
	return item, nil
}
