package repository_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/m9rc1n/shop/app/items-api/repository"
	"github.com/m9rc1n/shop/pkg/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRepository(t *testing.T) {
	logger := log.New()
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db := sqlx.NewDb(mockDB, "sqlmock")
	defer mockDB.Close()

	SUT := repository.NewRepository(db, logger)

	t.Run(
		"List Items - Success", func(t *testing.T) {
			// Given
			rows := sqlmock.NewRows([]string{"name", "quantity"}).
				AddRow("Item 1", 10).
				AddRow("Item 2", 20)
			mock.ExpectQuery(repository.ListQuery).WillReturnRows(rows)
			// When
			items, err := SUT.List()
			// Then
			assert.Nil(t, err)
			assert.Len(t, items, 2)
		},
	)

	t.Run(
		"Get Item - Success", func(t *testing.T) {
			// Given
			rows := sqlmock.NewRows([]string{"id", "name", "quantity"}).
				AddRow(2, "Item 2", 20)
			mock.ExpectQuery("SELECT id").WithArgs(2).WillReturnRows(rows)
			// When
			item, err := SUT.Get(2)
			// Then
			assert.Nil(t, err)
			assert.Equal(t, int32(20), item.Quantity)
			assert.Equal(t, "Item 2", item.Name)
		},
	)

	t.Run(
		"Reserve Item - Success", func(t *testing.T) {
			// Given
			mock.ExpectExec("UPDATE shop.items SET").
				WithArgs(2, 1).
				WillReturnResult(sqlmock.NewResult(1, 1))
			// When
			err := SUT.Reserve(1, 2)
			// Then
			assert.Nil(t, err)
		},
	)

	t.Run(
		"Create Item - Success", func(t *testing.T) {
			// Given
			rows := sqlmock.NewRows([]string{"id"}).
				AddRow(5)
			mock.ExpectQuery("INSERT INTO shop.items").
				WithArgs("Item", 10).
				WillReturnRows(rows)
			// When
			itemID, err := SUT.Create("Item", 10)
			// Then
			assert.Nil(t, err)
			assert.Equal(t, int64(5), itemID)
		},
	)

	t.Run(
		"Create Item - Fail", func(t *testing.T) {
			// Given
			rows := sqlmock.NewRows([]string{"id"}).
				AddRow(5)
			mock.ExpectQuery("INSERT INTO shop.items").
				WithArgs("Item", 10).
				WillReturnRows(rows)
			// When
			itemID, err := SUT.Create("Item", 2)
			// Then
			assert.Error(t, err)
			assert.Equal(t, int64(-1), itemID)
		},
	)
}
