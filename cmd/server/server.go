package main

import (
	"context"
	"fmt"
	"github.com/MaksimPerv/Metric/config/serverconfig"
	"github.com/MaksimPerv/Metric/internal/db"
	"github.com/MaksimPerv/Metric/internal/fileutils"
	"github.com/MaksimPerv/Metric/internal/handlers"
	"github.com/MaksimPerv/Metric/internal/middleware"
	"github.com/MaksimPerv/Metric/internal/mistake"
	"github.com/MaksimPerv/Metric/internal/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
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
	projectRoot, err := filepath.Abs("../..")
	if err != nil {
		panic(fmt.Errorf("failed to get project root: %v", err))
	}
	filePath := filepath.Join(projectRoot, "internal", "profiles", "base.pprof")
	fmem, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}

	err = Initalize()
	if err != nil {
		panic(err)
	}
	serverconfig.ParseFlags()

	go func() {
		http.ListenAndServe(serverconfig.FlagRunAddr, run())
	}()
	defer fileutils.File.Close()

	defer db.DB.Close()

	time.Sleep(5 * time.Second)
	runtime.GC()
	pprof.WriteHeapProfile(fmem)

	select {}

}

func run() chi.Router {
	storage := storage.NewMemStorage()
	router := chi.NewRouter()

	if serverconfig.FileStoragePath != "" {
		fileutils.File, _ = fileutils.NewProducer(serverconfig.FileStoragePath)
		//defer fileutils.File.Close()
		if serverconfig.Restore {
			storage.Restore(fileutils.File.Read())
		}
		if serverconfig.StoreInterval != 0 {
			log.Println(serverconfig.FileStoragePath)
			log.Println(serverconfig.StoreInterval)
			go func() {
				for {
					err := fileutils.File.Write(storage.GetList())
					if err != nil {
						log.Print(err)
					}
					time.Sleep(serverconfig.StoreInterval)
				}
			}()
		}
	}

	router.Use(middleware.GzipMiddleware)
	router.Use(middleware.SignatureMiddleware)

	err := mistake.Retry(3, time.Second, func() error {
		err := db.Init(serverconfig.DatabaseDSN)
		return err
	})
	if err == nil {
		router.Get("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
				return
			}
			err = db.DB.PingContext(context.Background())
			if err != nil {
				http.Error(w, "Not connect", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(200)
			w.Write([]byte("Connect OK"))
		}))
	}
	if err != nil {
		panic(err)
	}
	handler := hendlers.NewMetricsHandlers(storage)

	database := db.Postgres{Database: db.DB}
	if serverconfig.DatabaseDSN != "" {
		handler = hendlers.NewMetricsHandlers(database)
	}

	router.Route("/", func(r chi.Router) {

		r.Get("/", (handler.GetList))
		//r.Get("/ping", handler.Ping)

		r.Post("/updates/", handler.Updates)

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
