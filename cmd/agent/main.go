package main

import (
	"github.com/MaksimPerv/Metric/config"
	"github.com/MaksimPerv/Metric/internal/collector"
	"github.com/MaksimPerv/Metric/internal/sender"
	"time"
)

func main() {
	// Загрузка конфигурации
	cfg := config.LoadConfig()

	// Инициализация коллектора и отправителя
	metricsCollector := collector.NewMetricCollector()
	metricsSender := sender.NewMetricsSender(cfg.ServerAddress)

	// Запуск сбора метрик
	go func() {
		for {
			metricsCollector.Collect()
			time.Sleep(cfg.PollInterval)
		}
	}()
	time.Sleep(time.Second * 2)
	// Запуск отправки метрик
	go func() {
		for {
			metrics := metricsCollector.GetMetrics()
			metricsSender.Send(metrics)
			time.Sleep(cfg.ReportInterval)
		}
	}()

	// Бесконечный цикл для работы программы
	select {}
}
