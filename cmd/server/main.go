package main

import (
	"github.com/MaksimPerv/Metric/internal/handlers"
	"github.com/MaksimPerv/Metric/internal/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	http.ListenAndServe(":8080", run())
}

func run() chi.Router {

	storage := storage.NewMemStorage()
	router := chi.NewRouter()
	handler := hendlers.NewMetricsHandlers(storage)
	router.Route("/", func(r chi.Router) {

		r.Get("/", handler.GetList)

		r.Route("/update", func(r chi.Router) {

			r.Route("/{type}", func(r chi.Router) {

				r.Route("/{name}", func(r chi.Router) {

					r.Post("/{value}", handler.UpdateMetric)
				})
			})
		})

		r.Route("/value", func(r chi.Router) {
			r.Route("/{type}", func(r chi.Router) {

				r.Get("/{name}", handler.GetMetric)

			})

		})
	})
	return router
}
