package collector

import (
	"github.com/MaksimPerv/Metric/pkg/metric"
	"math/rand"
	"runtime"
	"sync"
)

type MetricCollector struct {
	metrics metric.Metrics
	mu      sync.Mutex
}

func NewMetricCollector() *MetricCollector {
	return &MetricCollector{
		metrics: metric.Metrics{},
	}
}

func (c *MetricCollector) Collect() {
	c.mu.Lock()
	defer c.mu.Unlock()

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	c.metrics.Alloc = float64(memStats.Alloc)
	c.metrics.BuckHashSys = float64(memStats.BuckHashSys)
	c.metrics.Frees = float64(memStats.Frees)
	c.metrics.GCCPUFraction = memStats.GCCPUFraction
	c.metrics.GCSys = float64(memStats.GCSys)
	c.metrics.HeapAlloc = float64(memStats.HeapAlloc)
	c.metrics.HeapIdle = float64(memStats.HeapIdle)
	c.metrics.HeapInuse = float64(memStats.HeapInuse)
	c.metrics.HeapObjects = float64(memStats.HeapObjects)
	c.metrics.HeapReleased = float64(memStats.HeapReleased)
	c.metrics.HeapSys = float64(memStats.HeapSys)
	c.metrics.LastGC = float64(memStats.LastGC)
	c.metrics.Lookups = float64(memStats.Lookups)
	c.metrics.MCacheInuse = float64(memStats.MCacheInuse)
	c.metrics.MCacheSys = float64(memStats.MCacheSys)
	c.metrics.MSpanInuse = float64(memStats.MSpanInuse)
	c.metrics.MSpanSys = float64(memStats.MSpanSys)
	c.metrics.Mallocs = float64(memStats.Mallocs)
	c.metrics.NextGC = float64(memStats.NextGC)
	c.metrics.NumForcedGC = float64(memStats.NumForcedGC)
	c.metrics.NumGC = float64(memStats.NumGC)
	c.metrics.OtherSys = float64(memStats.OtherSys)
	c.metrics.PauseTotalNs = float64(memStats.PauseTotalNs)
	c.metrics.StackInuse = float64(memStats.StackInuse)
	c.metrics.StackSys = float64(memStats.StackSys)
	c.metrics.Sys = float64(memStats.Sys)
	c.metrics.TotalAlloc = float64(memStats.TotalAlloc)

	// Обновление PollCount
	c.metrics.PollCount++

	// Генерация RandomValue
	c.metrics.RandomValue = rand.Float64()
}

func (c *MetricCollector) GetMetrics() metric.Metrics {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.metrics
}
