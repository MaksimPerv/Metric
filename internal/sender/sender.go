package sender

import (
	"github.com/MaksimPerv/Metric/pkg/metric"
	"log"
	"net/http"
	"strconv"
)

type MetricsSender struct {
	serverAddress string
}

func NewMetricsSender(serverAddress string) *MetricsSender {
	return &MetricsSender{
		serverAddress: serverAddress,
	}
}
func (s *MetricsSender) Send(m metric.Metrics) {

	// Отправка метрик типа gauge
	for name, value := range m.GaugeMetrics() {
		url := s.serverAddress + "/update/gauge/" + name + "/" + strconv.FormatFloat(value, 'f', -1, 64)
		s.sendRequest(url)
	}

	// Отправка метрик типа counter
	for name, value := range m.CounterMetrics() {
		url := s.serverAddress + "/update/counter/" + name + "/" + strconv.FormatInt(value, 10)
		s.sendRequest(url)
	}
}

func (s *MetricsSender) sendRequest(url string) {
	resp, err := http.Post(url, "text/plain", nil)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Server returned non-OK status: %d", resp.StatusCode)
	}
}
