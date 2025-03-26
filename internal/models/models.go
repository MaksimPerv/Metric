package models

type MetricType string

const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)

type Metric struct {
	Name  string
	Type  MetricType
	Value interface{}
}
