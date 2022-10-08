package store_test

import (
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/m9rc1n/shop/app/reservations-api/api"
	"github.com/m9rc1n/shop/app/reservations-api/store"
	"github.com/m9rc1n/shop/pkg/log"
	"net/http"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	middleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	"github.com/deepmap/oapi-codegen/pkg/testutil"
)

func TestStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var err error
	// Get the swagger description of our API
	swagger, err := api.GetSwagger()
	require.NoError(t, err)
	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil
	// This is how you set up a basic chi router
	r := chi.NewRouter()
	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	r.Use(middleware.OapiRequestValidator(swagger))
	SUT := store.NewStore(log.New(), time.Microsecond)
	api.HandlerFromMux(SUT, r)

	t.Run(
		"Create Reservation - Success", func(t *testing.T) {
			// Given
			request := api.CreateReservationJSONRequestBody{
				ItemName:     "ItemName",
				ItemQuantity: 1,
			}
			// When
			rr := testutil.NewRequest().Post("/reserve").WithJsonBody(request).GoWithHTTPHandler(t, r).Recorder
			// Then
			var result api.Reservation
			err = json.NewDecoder(rr.Body).Decode(&result)
			assert.NoError(t, err, "error unmarshalling response")
			assert.Equal(t, request.ItemName, result.ItemName)
			assert.Equal(t, request.ItemQuantity, result.ItemQuantity)
			assert.Equal(t, http.StatusCreated, rr.Code)
			assert.NotNil(t, result.Id)
		},
	)

	t.Run(
		"Is Alive", func(t *testing.T) {
			// Given
			// When
			rr := testutil.NewRequest().Get("/alive").GoWithHTTPHandler(t, r).Recorder
			// Then
			assert.Equal(t, http.StatusOK, rr.Code)
		},
	)

	t.Run(
		"Is Ready", func(t *testing.T) {
			// Given
			// When
			rr := testutil.NewRequest().Get("/ready").GoWithHTTPHandler(t, r).Recorder
			// Then
			assert.Equal(t, http.StatusOK, rr.Code)
		},
	)
}
