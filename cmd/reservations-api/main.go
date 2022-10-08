package main

import (
	"github.com/m9rc1n/shop/app/reservations-api/server"
	"github.com/m9rc1n/shop/app/reservations-api/store"
	"github.com/m9rc1n/shop/pkg/log"
	"time"
)

func main() {
	newLogger := log.New()
	newStore := store.NewStore(newLogger, time.Second)
	server.Run(newStore, newLogger)
}
