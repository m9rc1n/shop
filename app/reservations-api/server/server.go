package server

import (
	"fmt"
	middleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	"github.com/go-chi/chi/v5"
	oamiddleware "github.com/go-openapi/runtime/middleware"
	"github.com/m9rc1n/shop/app/reservations-api/api"
	"github.com/m9rc1n/shop/pkg/log"
	"net/http"
)

func Run(store api.ServerInterface, logger log.Logger) {
	swagger, err := api.GetSwagger()
	if err != nil {
		panic(fmt.Sprintf("Error loading swagger spec: %s", err))
	}
	swagger.Servers = nil
	r := chi.NewRouter()
	r.Use(
		func(next http.Handler) http.Handler {
			return oamiddleware.SwaggerUI(
				oamiddleware.SwaggerUIOpts{
					SpecURL: "static/swagger.yml",
				}, next,
			)
		},
	)
	r.Route(
		"/", func(r chi.Router) {
			r.Group(
				func(r chi.Router) {
					fileServer := http.FileServer(http.Dir("./static/"))
					r.Handle("/static/*", http.StripPrefix("/static", fileServer))
				},
			)
			r.Group(
				func(r chi.Router) {
					r.Use(middleware.OapiRequestValidator(swagger))
					api.HandlerFromMux(store, r)
				},
			)
		},
	)
	s := &http.Server{
		Handler: r,
		Addr:    ":8080",
	}
	logger.Info("Started Reservations API Server")
	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}
}
