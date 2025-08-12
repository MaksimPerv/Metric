package hendlers

import (
	"encoding/json"
	"fmt"
	"github.com/MaksimPerv/Metric/internal/models"
	"github.com/MaksimPerv/Metric/internal/storage"
	"github.com/go-chi/chi/v5"
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

func (h *MetricsHandlers) GetList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	result := h.storage.GetList()
	strData := fmt.Sprintf("%v", result)
	if strings.Contains(r.Header.Get("Accept"), "text/html") {
		w.Header().Set("Content-Type", "text/html")
	} else {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}
	w.Header().Set("Date", time.Now().Format(time.RFC1123))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strData))
}

func (h *MetricsHandlers) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	//if r.Header.Get("Content-Type") != "text/plain" {
	//
	//	http.Error(w, "Invalid Content-Type", http.StatusUnsupportedMediaType)
	//	return
	//}
	//if r.Header.Get("Content-Length") != "0" {
	//	http.Error(w, "Invalid Content-Length", http.StatusBadRequest)
	//	return
	//}
	//parts := strings.Split(r.URL.Path, "/")
	//if len(parts) != 5 {
	//	w.WriteHeader(http.StatusNotFound)
	//	w.Write([]byte("metrics not found"))
	//	return
	//}
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")
	metricValue := chi.URLParam(r, "value")
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
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	metricName := chi.URLParam(r, "name")
	//metricType := chi.URLParam(r, "type")

	result, ok := h.storage.GetMetric(metricName)
	if !ok {
		http.Error(w, "Not Found Metric", http.StatusNotFound)
		return
	}
	//w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	//w.Header().Set("Date", time.Now().Format(time.RFC1123))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(getValueAsString(result.Value)))
}

func getValueAsString(value interface{}) string {
	switch v := value.(type) {
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		return ""
	}
}

func (h *MetricsHandlers) UpdateJSONMetric(w http.ResponseWriter, r *http.Request) {
	var req models.Metrics
	if ct := r.Header.Get("Content-Type"); ct != "application/json" {
		http.Error(w, "Invalid content type", http.StatusBadRequest)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Problem Body", http.StatusBadRequest)
		return
	}
	//log.Println(req)
	var rawMetric models.Metric
	switch req.MType {
	case string(models.Gauge):
		rawMetric = models.Metric{
			Name:  req.ID,
			Type:  models.Gauge,
			Value: *req.Value,
		}
	case string(models.Counter):
		rawMetric = models.Metric{
			Name:  req.ID,
			Type:  models.Counter,
			Value: *req.Delta,
		}

	}
	rawMetric = h.storage.UpdateMetric(rawMetric)

	switch rawMetric.Type {
	case models.Gauge:
		*req.Value = rawMetric.Value.(float64)
	case models.Counter:
		*req.Delta = rawMetric.Value.(int64)
	}
	response, _ := json.Marshal(req)

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	//log.Println(string(response))
	w.Write(response)
}

func (h *MetricsHandlers) GetJSONMetric(w http.ResponseWriter, r *http.Request) {
	if ct := r.Header.Get("Content-Type"); ct != "" && ct != "application/json" {
		http.Error(w, "Invalid content type", http.StatusBadRequest)
		return
	}
	var req models.Metrics
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Not JSOn", http.StatusBadRequest)
		return
	}
	metric, ok := h.storage.GetMetric(req.ID)
	if !ok {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	switch metric.Type {
	case models.Gauge:

		val, _ := metric.Value.(float64)
		req.Value = &val
	case models.Counter:
		val, _ := metric.Value.(int64)
		req.Delta = &val
	}
	response, _ := json.Marshal(req)
	//log.Println(req)
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (h *MetricsHandlers) Updates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var req []models.Metrics
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Can't decode json", http.StatusBadRequest)
		return
	}
	if len(req) == 0 {
		http.Error(w, "Empty metrics array", http.StatusBadRequest)
		return
	}
	for _, value := range req {
		var rawMetric models.Metric
		switch value.MType {
		case string(models.Gauge):
			if value.Value == nil {
				http.Error(w, fmt.Sprintf("Value is required for gauge metric %s", value.ID),
					http.StatusBadRequest)
				return
			}
			rawMetric = models.Metric{
				Name:  value.ID,
				Type:  models.Gauge,
				Value: *value.Value,
			}
		case string(models.Counter):
			if value.Delta == nil {
				http.Error(w, fmt.Sprintf("Delta is required for counter metric %s", value.ID),
					http.StatusBadRequest)
				return
			}
			rawMetric = models.Metric{
				Name:  value.ID,
				Type:  models.Counter,
				Value: *value.Delta,
			}

		}
		h.storage.UpdateMetric(rawMetric)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Metric update"))

}

//func (h *MetricsHandlers) Ping(w http.ResponseWriter, r *http.Request) {
//	if r.Method != http.MethodGet {
//		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
//		return
//	}
//	err := db.DB.PingContext(context.Background())
//	if err != nil {
//		http.Error(w, "Not connect", http.StatusInternalServerError)
//		return
//	}
//	w.WriteHeader(200)
//	w.Write([]byte("Connect OK"))
//}
