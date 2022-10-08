package store_test

import (
	"encoding/json"
	"github.com/aws/smithy-go/ptr"
	"github.com/golang/mock/gomock"
	"github.com/m9rc1n/shop/app/items-api/api"
	mockservice "github.com/m9rc1n/shop/app/items-api/service/mock"
	"github.com/m9rc1n/shop/app/items-api/store"
	"github.com/m9rc1n/shop/pkg/log"
	"github.com/pkg/errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	middleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	"github.com/deepmap/oapi-codegen/pkg/testutil"
)

func doGet(t *testing.T, mux *chi.Mux, url string) *httptest.ResponseRecorder {
	response := testutil.NewRequest().Get(url).WithAcceptJson().GoWithHTTPHandler(t, mux)
	return response.Recorder
}

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
	serviceMock := mockservice.NewMockService(ctrl)
	SUT := store.NewStore(log.New(), serviceMock)
	api.HandlerFromMux(SUT, r)

	t.Run(
		"Create Item - Success", func(t *testing.T) {
			// Given
			request := api.CreateItemJSONRequestBody{
				Name:     "ItemName",
				Quantity: 1,
			}
			serviceMock.EXPECT().
				CreateItem(gomock.Any()).
				Return(
					&api.Item{
						Id:       ptr.Int64(1),
						Quantity: request.Quantity,
						Name:     request.Name,
					}, nil,
				).
				Times(1)
			// When
			rr := testutil.NewRequest().Post("/items").WithJsonBody(request).GoWithHTTPHandler(t, r).Recorder
			// Then
			var result api.Item
			err = json.NewDecoder(rr.Body).Decode(&result)
			assert.NoError(t, err, "error unmarshalling response")
			assert.Equal(t, request.Name, result.Name)
			assert.Equal(t, request.Quantity, result.Quantity)
			assert.Equal(t, http.StatusCreated, rr.Code)
			assert.NotNil(t, result.Id)
		},
	)
	t.Run(
		"Create Item - Corrupted Request", func(t *testing.T) {
			// Given
			serviceMock.EXPECT().
				CreateItem(gomock.Any()).
				Times(0)
			// When
			rr := testutil.NewRequest().Post("/items").WithJsonBody("<>").GoWithHTTPHandler(t, r).Recorder
			// Then
			assert.Equal(t, http.StatusBadRequest, rr.Code)
		},
	)
	t.Run(
		"Create Item - Corrupted Response", func(t *testing.T) {
			// Given
			request := api.CreateItemJSONRequestBody{
				Name:     "ItemName",
				Quantity: 1,
			}
			serviceMock.EXPECT().
				CreateItem(gomock.Any()).
				Return(nil, errors.New("error")).
				Times(1)
			// When
			rr := testutil.NewRequest().Post("/items").WithJsonBody(request).GoWithHTTPHandler(t, r).Recorder
			// Then
			var result api.Error
			err = json.NewDecoder(rr.Body).Decode(&result)
			assert.Equal(t, http.StatusBadRequest, rr.Code)
			assert.Equal(t, http.StatusBadRequest, result.Code)
			assert.Equal(t, "Cannot store the item: error", result.Message)
		},
	)
	t.Run(
		"List Items - Success", func(t *testing.T) {
			// Given
			serviceMock.EXPECT().
				ListItems().
				Return(
					[]*api.Item{
						{
							Id:       ptr.Int64(1),
							Quantity: 10,
							Name:     "Item 1",
						},
						{
							Id:       ptr.Int64(2),
							Quantity: 10,
							Name:     "Item 2",
						},
					}, nil,
				).
				Times(1)
			// When
			rr := doGet(t, r, "/items")
			// Then
			var items []*api.Item
			err = json.NewDecoder(rr.Body).Decode(&items)
			assert.Equal(t, http.StatusOK, rr.Code)
			assert.NoError(t, err, "error getting response", err)
			assert.Equal(t, 2, len(items))
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
