package models

import (
	"encoding/json"
	"fmt"
)

type MetricType string

const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)

type Metric struct {
	Name  string      `json:"name"`
	Type  MetricType  `json:"type"`
	Value interface{} `json:"value"`
}

func (m *Metric) UnmarshalJSON(data []byte) error {
	// Временная структура для парсинга
	var temp struct {
		Name  string          `json:"name"`
		Type  MetricType      `json:"type"`
		Value json.RawMessage `json:"value"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	m.Name = temp.Name
	m.Type = temp.Type

	// Обрабатываем Value в зависимости от типа метрики
	switch m.Type {
	case Gauge:
		var v float64
		if err := json.Unmarshal(temp.Value, &v); err != nil {
			return err
		}
		m.Value = v
	case Counter:
		var v int64
		if err := json.Unmarshal(temp.Value, &v); err != nil {
			return err
		}
		m.Value = v
	default:
		return fmt.Errorf("unknown metric type: %s", m.Type)
	}

	return nil
}

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}
