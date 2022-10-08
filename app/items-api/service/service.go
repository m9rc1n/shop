package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/m9rc1n/shop/app/items-api/api"
	"github.com/m9rc1n/shop/app/items-api/repository"
	rapi "github.com/m9rc1n/shop/app/reservations-api/api"
	"github.com/m9rc1n/shop/pkg/log"
	"github.com/pkg/errors"
)

// Service define item service methods interface
type Service interface {
	// ListItems returns the list of items
	ListItems() ([]*api.Item, error)
	// CreateItem based on rest API params and persist it
	CreateItem(item *api.Item) (*api.Item, error)
}

// service persists item in database
type service struct {
	repository         repository.Repository
	logger             log.Logger
	reservationsClient rapi.ClientInterface
}

// NewService creates a new item service
func NewService(repository repository.Repository, reservationsClient rapi.ClientInterface, log log.Logger) Service {
	return service{
		reservationsClient: reservationsClient,
		repository:         repository,
		logger:             log,
	}
}

// ListItems based on request parameters
func (r service) ListItems() ([]*api.Item, error) {
	storageItems, err := r.repository.List()
	if err != nil {
		return nil, errors.Wrap(err, "error fetching items from storage")
	}
	apiItems := make([]*api.Item, 0)
	for _, storageItem := range storageItems {
		apiItems = append(apiItems, mapStorageItemToApiItem(storageItem))
	}
	return apiItems, nil
}

// CreateItem from name and quantity
func (r service) CreateItem(item *api.Item) (*api.Item, error) {
	itemId, err := r.repository.Create(item.Name, item.Quantity)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("error creating item in storage: %v", item))
	}
	storedItem, err := r.repository.Get(itemId)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("error fetching item from storage: %v", item))
	}
	// Run goroutine in order to reserve the item
	go r.reserveItem(storedItem)
	return mapStorageItemToApiItem(storedItem), nil
}

func (r service) reserveItem(storedItem *repository.Item) {
	createReservationResponse, err := r.reservationsClient.CreateReservation(
		context.Background(), rapi.CreateReservationJSONBody{
			ItemQuantity: storedItem.Quantity,
			ItemName:     storedItem.Name,
		},
	)
	if err != nil {
		r.logger.Error(errors.Wrap(err, fmt.Sprintf("error creating reservation for item id: %d", *storedItem.ID)))
		return
	}
	newReservation := new(rapi.Reservation)
	if err := json.NewDecoder(createReservationResponse.Body).Decode(newReservation); err != nil {
		r.logger.Error(errors.Wrap(err, fmt.Sprintf("invalid format for the reservation")))
		return
	}
	if err := r.repository.Reserve(*storedItem.ID, *newReservation.Id); err != nil {
		r.logger.Error(errors.Wrap(err, fmt.Sprintf("error persisting reservation for item id: %d", *storedItem.ID)))
	}
}

func mapStorageItemToApiItem(storageItem *repository.Item) *api.Item {
	return &api.Item{
		Id:            storageItem.ID,
		CreatedAt:     storageItem.CreatedAt,
		UpdatedAt:     storageItem.UpdatedAt,
		ReservedAt:    storageItem.ReservedAt,
		ReservationId: storageItem.ReservationID,
		Name:          storageItem.Name,
		Quantity:      storageItem.Quantity,
	}
}
