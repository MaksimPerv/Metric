package main

import (
	"github.com/MaksimPerv/Metric/config/serverconfig"
	"github.com/MaksimPerv/Metric/internal/handlers"
	"github.com/MaksimPerv/Metric/internal/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type (

	// берём структуру для хранения сведений об ответе
	responseData struct {
		status int
		size   int
	}

	// добавляем реализацию http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

var Log *zap.Logger = zap.NewNop()

func Initalize() error {

	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)

	zl, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	Log = zl
	return nil
}

func RequestLogger(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()
		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: writer, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}
		next.ServeHTTP(&lw, request)
		Log.Info("got incoming HHTP request",
			zap.String("method", request.Method),
			zap.String("url", request.URL.Path),
			zap.Duration("time", time.Since(start)),
			zap.Int("status", responseData.status),
			zap.Int("size", responseData.size),
		)
	}
}
func main() {
	err := Initalize()
	if err != nil {
		panic(err)
	}

	serverconfig.ParseFlags()
	http.ListenAndServe(serverconfig.FlagRunAddr, run())
}

func run() chi.Router {
	storage := storage.NewMemStorage()
	router := chi.NewRouter()
	handler := hendlers.NewMetricsHandlers(storage)
	//router.Use(middleware.GzipMiddleware)
	router.Route("/", func(r chi.Router) {

		r.Get("/", (handler.GetList))

		r.Route("/update", func(r chi.Router) {

			r.Post("/", (handler.UpdateJSONMetric))

			r.Route("/{type}", func(r chi.Router) {

				r.Route("/{name}", func(r chi.Router) {

					r.Post("/{value}", (handler.UpdateMetric))
				})
			})
		})

		r.Route("/value", func(r chi.Router) {
			r.Post("/", (handler.GetJSONMetric))
			r.Route("/{type}", func(r chi.Router) {

				r.Get("/{name}", (handler.GetMetric))

			})

		})
	})

	return router
}
