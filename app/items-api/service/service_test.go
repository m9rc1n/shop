package service_test

import (
	"bytes"
	"github.com/aws/smithy-go/ptr"
	"github.com/golang/mock/gomock"
	"github.com/m9rc1n/shop/app/items-api/api"
	"github.com/m9rc1n/shop/app/items-api/repository"
	mockrepository "github.com/m9rc1n/shop/app/items-api/repository/mock"
	"github.com/m9rc1n/shop/app/items-api/service"
	mockapi "github.com/m9rc1n/shop/app/reservations-api/api/mock"
	"github.com/m9rc1n/shop/pkg/log"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repositoryMock := mockrepository.NewMockRepository(ctrl)
	apiReservationsClientMock := mockapi.NewMockClientInterface(ctrl)
	SUT := service.NewService(repositoryMock, apiReservationsClientMock, log.New())

	t.Run(
		"Create Item - Success", func(t *testing.T) {
			// Given
			item := &api.Item{
				Name:     "ItemName",
				Quantity: 1,
			}
			itemID := int64(2)
			storedItem := &repository.Item{
				ID:       ptr.Int64(itemID),
				Name:     "ItemName",
				Quantity: 1,
			}
			repositoryMock.EXPECT().
				Create(item.Name, item.Quantity).
				Return(itemID, nil).
				Times(1)
			repositoryMock.EXPECT().
				Get(itemID).
				Return(storedItem, nil).
				Times(1)
			reservationJSON := `{"id":123,"item_name":"ItemName","item_quantity":1}`
			// create a new reader with that JSON
			reservationResponse := ioutil.NopCloser(bytes.NewReader([]byte(reservationJSON)))
			apiReservationsClientMock.EXPECT().
				CreateReservation(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(
					&http.Response{Body: reservationResponse, StatusCode: http.StatusCreated}, nil,
				).
				Times(1)
			repositoryMock.EXPECT().
				Reserve(itemID, int64(123)).
				Return(nil).
				Times(1)
			// When
			newItem, err := SUT.CreateItem(item)
			time.Sleep(time.Millisecond * 100)
			// Then
			assert.Nil(t, err)
			assert.NotNil(t, newItem.Id)
			assert.Equal(t, storedItem.ID, newItem.Id)
		},
	)

	t.Run(
		"List Items - Success", func(t *testing.T) {
			// Given
			repositoryMock.EXPECT().
				List().
				Return(
					[]*repository.Item{
						{
							ID:       ptr.Int64(1),
							Name:     "ItemName 1",
							Quantity: 1,
						},
						{
							ID:       ptr.Int64(2),
							Name:     "ItemName 2",
							Quantity: 5,
						},
					}, nil,
				).
				Times(1)
			// When
			actualItems, err := SUT.ListItems()
			// Then
			assert.Nil(t, err)
			assert.Len(t, actualItems, 2)
			assert.Equal(t, actualItems[1].Name, "ItemName 2")
			assert.Equal(t, actualItems[0].Quantity, int32(1))
		},
	)
}
