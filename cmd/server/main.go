package main

import (
	"github.com/MaksimPerv/Metric/internal/handlers"
	"github.com/MaksimPerv/Metric/internal/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	run()
}

func run() error {

	storage := storage.NewMemStorage()
	router := chi.NewRouter()
	handler := hendlers.NewMetricsHandlers(storage)
	router.Route("/update", func(r chi.Router) {

		r.Route("/{type}", func(r chi.Router) {

			r.Route("/{name}", func(r chi.Router) {

				r.Post("/{value}", handler.UpdateMetric)
			})
		})
	})

	router.Route("/value", func(r chi.Router) {
		r.Route("/{type}", func(r chi.Router) {

			r.Get("/{name}", handler.GetMetric)

		})

	})

	if err := http.ListenAndServe(":8080", router); err != nil {
		return err
	}
	return nil
}
