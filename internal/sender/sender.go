package sender

import (
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
