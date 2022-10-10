package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/m9rc1n/shop/app/items-api/repository"
	"github.com/m9rc1n/shop/app/items-api/server"
	"github.com/m9rc1n/shop/app/items-api/service"
	"github.com/m9rc1n/shop/app/items-api/store"
	rapi "github.com/m9rc1n/shop/app/reservations-api/api"
	"github.com/m9rc1n/shop/internal/config"
	"github.com/m9rc1n/shop/pkg/log"
)

func main() {
	newLogger := log.New()
	newConfig, err := config.New(newLogger)
	if err != nil {
		newLogger.Errorf("failed to initialize configuration: %s", err.Error())
		panic(err)
	}
	newDB := sqlx.MustConnect("postgres", newConfig.DB.Dsn)
	newReservationsClient, err := rapi.NewClient(newConfig.ReservationsEndpoint)
	if err != nil {
		newLogger.Errorf("failed to initialize reservations client: %s", err.Error())
		panic(err)
	}
	newRepository := repository.NewRepository(newDB, newLogger)
	newMigrator := repository.NewMigrator(newDB, newLogger)
	newService := service.NewService(
		newRepository,
		newReservationsClient,
		newLogger,
	)
	newStore := store.NewStore(newLogger, newService)
	if err := newMigrator.MigrateUp(); err != nil {
		newLogger.Errorf("failed to migrating items: %s", err.Error())
		panic(err)
	}
	server.Run(newStore, newLogger)
}
