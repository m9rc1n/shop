package store

import (
	"encoding/json"
	"fmt"
	"github.com/m9rc1n/shop/app/items-api/api"
	"github.com/m9rc1n/shop/app/items-api/service"
	"github.com/m9rc1n/shop/pkg/log"
	"net/http"
)

type store struct {
	service service.Service
	logger  log.Logger
}

func NewStore(logger log.Logger, service service.Service) api.ServerInterface {
	return &store{
		logger:  logger,
		service: service,
	}
}

func (s store) ListItems(w http.ResponseWriter, _ *http.Request) {
	storedItems, err := s.service.ListItems()
	if err != nil {
		s.sendError(w, http.StatusBadRequest, "Cannot fetch list of the items")
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(storedItems); err != nil {
		s.logger.Error(err)
	}
}

func (s store) CreateItem(w http.ResponseWriter, r *http.Request) {
	newItem := new(api.Item)
	if err := json.NewDecoder(r.Body).Decode(newItem); err != nil {
		s.sendError(w, http.StatusBadRequest, "Invalid format for the item")
		return
	}
	storedItem, err := s.service.CreateItem(newItem)
	if err != nil {
		s.sendError(w, http.StatusBadRequest, fmt.Sprintf("Cannot store the item: %s", err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(storedItem); err != nil {
		s.logger.Error(err)
	}
}

func (s store) IsAlive(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(200)
}

func (s store) IsReady(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(200)
}

func (s store) sendError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(
		&api.Error{
			Code:    code,
			Message: message,
		},
	); err != nil {
		s.logger.Error(err)
	}
}
