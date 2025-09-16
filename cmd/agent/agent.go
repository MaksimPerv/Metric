package main

import (
	"github.com/MaksimPerv/Metric/config/agentconfig"
	"github.com/MaksimPerv/Metric/internal/collector"
	"github.com/MaksimPerv/Metric/internal/mistake"
	"github.com/MaksimPerv/Metric/internal/sender"
	"github.com/go-resty/resty/v2"
	"strconv"
	"sync"
	"time"
)

func main() {
	agentconfig.ParseFlags()
	// Загрузка конфигурации
	client := resty.New()
	var wg sync.WaitGroup
	// Инициализация коллектора и отправителя
	metricsCollector := collector.NewMetricCollector()
	metricsSender := sender.NewMetricsSender(client, agentconfig.FlagRunAddr)
	rateLimit, _ := strconv.Atoi(agentconfig.RateLimit)
	// Запуск сбора метрик
	go func() {
		for {
			metricsCollector.Collect()
			time.Sleep(agentconfig.PollInterval)
		}
	}()
	time.Sleep(1 * time.Second)
	// Запуск отправки метрик

	for i := 0; i < rateLimit; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				metrics := metricsCollector.GetMetrics()
				//metricsSender.Send(metrics)
				//metricsSender.SendJSON(metrics)
				err := mistake.Retry(3, time.Second, func() error {
					a := metricsSender.SendsMetrics(metrics)
					return a
				})
				if err != nil {
					panic(err)
				}
				time.Sleep(agentconfig.ReportInterval)
			}
		}()
	}
	go func() {
		batch := collector.GetMemMetric()
		err := mistake.Retry(3, time.Second, func() error {
			a := metricsSender.SendMemMetricsBatch(batch)
			return a
		})
		if err != nil {
			panic(err)
		}
		time.Sleep(agentconfig.ReportInterval)
	}()

	// Бесконечный цикл для работы программы
	wg.Wait()
}
