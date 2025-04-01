package sender

import (
	"encoding/json"
	"github.com/MaksimPerv/Metric/internal/models"
	"github.com/MaksimPerv/Metric/pkg/metric"
	"github.com/go-resty/resty/v2"
	"log"
	"net/http"
	"strconv"
)

type MetricsSender struct {
	serverAddress string
	Client        *resty.Client
}

func NewMetricsSender(client *resty.Client, serverAddress string) *MetricsSender {
	return &MetricsSender{
		serverAddress: serverAddress,
		Client:        client,
	}
}

func (s *MetricsSender) SendJson(m metric.Metrics) {
	for name, value := range m.GaugeMetrics() {
		url := "http://" + s.serverAddress + "/update"
		metrics := models.Metrics{
			ID:    name,
			MType: "gauge",
			Delta: nil,
			Value: &value,
		}
		obj, _ := json.Marshal(metrics)
		s.sendJsonRequest(url, obj)
	}
	for name, value := range m.CounterMetrics() {
		url := "http://" + s.serverAddress + "/update"
		metrics := models.Metrics{
			ID:    name,
			MType: "counter",
			Delta: &value,
			Value: nil,
		}
		obj, _ := json.Marshal(metrics)
		s.sendJsonRequest(url, obj)
	}
}

func (s *MetricsSender) sendJsonRequest(url string, value []byte) {
	var result models.Metrics
	resp, err := s.Client.R().SetHeader("Content-Type", "application/json").SetBody(value).SetResult(&result).Post(url)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return
	}
	log.Println(result)
	if resp.StatusCode() != http.StatusOK {
		log.Printf("Server returned non-OK status: %d", resp.StatusCode())

	}
}

func (s *MetricsSender) Send(m metric.Metrics) {
	// Отправка метрик типа gauge
	for name, value := range m.GaugeMetrics() {
		url := "http://" + s.serverAddress + "/update/gauge/" + name + "/" + strconv.FormatFloat(value, 'f', -1, 64)
		s.sendRequest(url)
	}

	// Отправка метрик типа counter
	for name, value := range m.CounterMetrics() {
		url := "http://" + s.serverAddress + "/update/counter/" + name + "/" + strconv.FormatInt(value, 10)
		s.sendRequest(url)
	}
}

func (s *MetricsSender) sendRequest(url string) {
	resp, err := s.Client.R().SetHeader("Content-Type", "text/plain").Post(url)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return
	}
	if resp.StatusCode() != http.StatusOK {
		log.Printf("Server returned non-OK status: %d", resp.StatusCode())

	}
}
