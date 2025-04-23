package main

import (
	"github.com/MaksimPerv/Metric/config/agentconfig"
	"github.com/MaksimPerv/Metric/internal/collector"
	"github.com/MaksimPerv/Metric/internal/sender"
	"github.com/go-resty/resty/v2"
	"time"
)

func main() {
	agentconfig.ParseFlags()
	// Загрузка конфигурации
	client := resty.New()
	// Инициализация коллектора и отправителя
	metricsCollector := collector.NewMetricCollector()
	metricsSender := sender.NewMetricsSender(client, agentconfig.FlagRunAddr)

	// Запуск сбора метрик
	go func() {
		for {
			metricsCollector.Collect()
			time.Sleep(agentconfig.PollInterval)
		}
	}()
	time.Sleep(time.Second * 2)
	// Запуск отправки метрик
	go func() {
		for {
			metrics := metricsCollector.GetMetrics()
			//metricsSender.Send(metrics)
			//metricsSender.SendJSON(metrics)
			metricsSender.SendsMetrics(metrics)
			time.Sleep(agentconfig.ReportInterval)
		}
	}()

	// Бесконечный цикл для работы программы
	select {}
}
