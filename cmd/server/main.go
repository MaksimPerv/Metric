package main

import (
	"github.com/MaksimPerv/Metric/internal/handlers"
	"github.com/MaksimPerv/Metric/internal/storage"
	"net/http"
)

func main() {
	run()
}

func run() error {

	storage := storage.NewMemStorage()

	handler := hendlers.NewMetricsHandlers(storage)

	http.HandleFunc("/", handler.UpdateMetric)
	//http.HandleFunc("/get", handler.GetMetrics)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		return err
	}
	return nil
}
