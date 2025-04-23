package sender

import (
	"bytes"
	"compress/gzip"
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

func (s *MetricsSender) SendJSON(m metric.Metrics) {
	for name, value := range m.GaugeMetrics() {
		url := "http://" + s.serverAddress + "/update"
		metrics := models.Metrics{
			ID:    name,
			MType: "gauge",
			Delta: nil,
			Value: &value,
		}
		obj, _ := json.Marshal(metrics)
		s.sendJSONRequest(url, obj)
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
		s.sendJSONRequest(url, obj)
	}
}

func (s *MetricsSender) sendJSONRequest(url string, value []byte) {
	var buf bytes.Buffer
	zp := gzip.NewWriter(&buf)
	zp.Write(value)
	zp.Close()
	//log.Println(string(buf.Bytes()))
	resp, err := s.Client.R().SetHeader("Content-Type", "application/json").SetHeader("Content-Encoding", "gzip").SetHeader("Accept-Encoding", "gzip").SetBody(buf.Bytes()).Post(url)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return
	}
	log.Println(string(resp.Body()))
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

//resp := s.Client.R().SetHeader("Content-Type", "application/json")
//resp.Header.Set("Content-Encoding", "gzip")
//var buf bytes.Buffer
//gz := gzip.NewWriter(&buf)
//_, err := gz.Write(value)
//if err != nil {
//panic(err)
//}
//gz.Close()
//resp.SetBody(buf.Bytes())
//response, err := resp.Post(url)
//if err != nil {
//log.Printf("Error sending request: %v", err)
//return
//}
//log.Println(response)
//if response.StatusCode() != http.StatusOK {
//log.Printf("Server returned non-OK status: %d", response.StatusCode())
//
//}

func (s *MetricsSender) SendsMetrics(m metric.Metrics) {
	delta := m.CounterMetrics()
	url := "http://" + s.serverAddress + "/updates/"
	value := m.GaugeMetrics()
	var batch []models.Metrics
	for name, v := range delta {
		met := models.Metrics{
			ID:    name,
			MType: "counter",
			Delta: &v,
			Value: nil,
		}
		batch = append(batch, met)
	}
	for name, v := range value {
		met := models.Metrics{
			ID:    name,
			MType: "gauge",
			Delta: nil,
			Value: &v,
		}
		batch = append(batch, met)
	}
	s.sendMetricsBatch(url, batch)

}
func (s *MetricsSender) sendMetricsBatch(url string, metrics []models.Metrics) {
	var buf bytes.Buffer
	obj, _ := json.Marshal(metrics)
	zp := gzip.NewWriter(&buf)
	zp.Write(obj)
	zp.Close()
	resp, err := s.Client.R().SetHeader("Content-Type", "application/json").SetHeader("Content-Encoding", "gzip").SetHeader("Accept-Encoding", "gzip").SetBody(buf.Bytes()).Post(url)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return
	}
	log.Println(string(resp.Body()))
	if resp.StatusCode() != http.StatusOK {
		log.Printf("Server returned non-OK status: %d", resp.StatusCode())

	}
}
