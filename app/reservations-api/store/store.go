package store

import (
	"encoding/json"
	"github.com/aws/smithy-go/ptr"
	"github.com/m9rc1n/shop/app/reservations-api/api"
	"github.com/m9rc1n/shop/pkg/log"
	"math/rand"
	"net/http"
	"time"
)

type store struct {
	logger      log.Logger
	sleeperUnit time.Duration
}

func NewStore(logger log.Logger, sleeperUnit time.Duration) api.ServerInterface {
	return &store{
		logger:      logger,
		sleeperUnit: sleeperUnit,
	}
}

func (s store) IsAlive(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(200)
}

func (s store) IsReady(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(200)
}

func (s store) CreateReservation(w http.ResponseWriter, r *http.Request) {
	time.Sleep(s.sleeperUnit * time.Duration(rand.Int63n(30)))
	newReservation := new(api.Reservation)
	if err := json.NewDecoder(r.Body).Decode(newReservation); err != nil {
		s.sendError(w, http.StatusBadRequest, "Invalid format for reservation")
		return
	}
	newReservation.Id = ptr.Int64(rand.Int63())
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newReservation); err != nil {
		s.logger.Error(err)
	}
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
