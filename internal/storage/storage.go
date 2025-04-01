package storage

import "github.com/MaksimPerv/Metric/internal/models"

type Storage interface {
	UpdateMetric(metric models.Metric) models.Metric
	GetMetric(name string) (models.Metric, bool)
	GetList() map[string]models.Metric
}

type MemStorage struct {
	metrics map[string]models.Metric
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		metrics: make(map[string]models.Metric),
	}
}

func (s *MemStorage) GetList() map[string]models.Metric {
	return s.metrics
}

func (s *MemStorage) UpdateMetric(metric models.Metric) models.Metric {
	if metric.Type == models.Counter {
		if existing, ok := s.metrics[metric.Name]; ok {
			metric.Value = existing.Value.(int64) + metric.Value.(int64)
		}
	}
	s.metrics[metric.Name] = metric
	return metric
}

func (s *MemStorage) GetMetric(name string) (models.Metric, bool) {
	metric, ok := s.metrics[name]
	return metric, ok
}
