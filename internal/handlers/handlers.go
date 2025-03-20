package hendlers

import (
	"github.com/MaksimPerv/Metric/internal/models"
	"github.com/MaksimPerv/Metric/internal/storage"

	"net/http"
	"strconv"
	"strings"
	"time"
)

type MetricsHandlers struct {
	storage storage.Storage
}

func NewMetricsHandlers(storage storage.Storage) *MetricsHandlers {
	return &MetricsHandlers{
		storage: storage,
	}
}

func (h *MetricsHandlers) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Invalid Content-Type", http.StatusUnsupportedMediaType)
		return
	}
	if r.Header.Get("Content-Length") != "0" {
		http.Error(w, "Invalid Content-Length", http.StatusBadRequest)
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 5 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("metrics not found"))
		return
	}
	metricType := parts[2]
	metricName := parts[3]
	metricValue := parts[4]
	var metric models.Metric
	metric.Name = metricName
	switch metricType {
	case string(models.Gauge):
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "Invalid gauge value", http.StatusBadRequest)
			return
		}
		metric.Type = models.Gauge
		metric.Value = value
	case string(models.Counter):
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "Invalid counter value", http.StatusBadRequest)
			return
		}
		metric.Type = models.Counter
		metric.Value = value
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("this type not found"))
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len("Metric updated\n")))
	w.Header().Set("Date", time.Now().Format(time.RFC1123))
	//log.Println(metric.Type, metric.Value)
	h.storage.UpdateMetric(metric)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Metric updated\n"))
}

func (h *MetricsHandlers) GetMetric(w http.ResponseWriter, r *http.Request) {

}
