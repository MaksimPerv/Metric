package db

import (
	"context"
	"database/sql"
	"github.com/MaksimPerv/Metric/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"strconv"
)

var DB *sql.DB

type Postgres struct {
	Database *sql.DB
}

func (p Postgres) UpdateMetric(metric models.Metric) models.Metric {
	// Преобразуем значение метрики в строку в зависимости от типа
	var valueStr string
	switch metric.Type {
	case models.Gauge:
		floatVal, ok := metric.Value.(float64)
		if !ok {
			return models.Metric{}
		}
		valueStr = strconv.FormatFloat(floatVal, 'f', -1, 64)
	case models.Counter:
		intVal, ok := metric.Value.(int64)
		if !ok {
			return models.Metric{}
		}
		valueStr = strconv.FormatInt(intVal, 10)
	default:
		return models.Metric{}
	}

	// Проверяем, существует ли метрика в базе
	var existingID int
	err := p.Database.QueryRow(
		"SELECT id FROM metrics WHERE name = $1 AND type = $2",
		metric.Name,
		string(metric.Type),
	).Scan(&existingID)

	if err != nil && err != sql.ErrNoRows {
		return models.Metric{}
	}

	if err == sql.ErrNoRows {
		// Метрика не существует, вставляем новую
		_, err = p.Database.Exec(
			"INSERT INTO metrics (name, type, value) VALUES ($1, $2, $3)",
			metric.Name,
			string(metric.Type),
			valueStr,
		)
		if err != nil {
			return models.Metric{}
		}
	} else {
		// Метрика существует, обновляем значение
		if metric.Type == models.Counter {
			// Для счетчика нужно добавить новое значение к существующему
			var currentValueStr string
			err = p.Database.QueryRow(
				"SELECT value FROM metrics WHERE id = $1",
				existingID,
			).Scan(&currentValueStr)
			if err != nil {
				return models.Metric{}
			}

			currentValue, err := strconv.ParseInt(currentValueStr, 10, 64)
			if err != nil {
				return models.Metric{}
			}

			newValue := currentValue + metric.Value.(int64)
			valueStr = strconv.FormatInt(newValue, 10)
			metric.Value = newValue
		}

		_, err = p.Database.Exec(
			"UPDATE metrics SET value = $1 WHERE id = $2",
			valueStr,
			existingID,
		)
		if err != nil {
			return models.Metric{}
		}
	}

	return metric
}

func (p Postgres) GetMetric(name string) (models.Metric, bool) {
	var typeMetric string
	var valueMEtric string
	log.Print(name)
	err := p.Database.QueryRow(`SELECT type,value FROM metrics WHERE name=$1`, name).Scan(&typeMetric, &valueMEtric)
	if err != nil && err != sql.ErrNoRows {
		return models.Metric{}, false
	}
	var result models.Metric
	result.Name = name
	if typeMetric == string(models.Counter) {
		result.Type = models.Counter
		currentValue, _ := strconv.ParseInt(valueMEtric, 10, 64)
		result.Value = currentValue
	} else {
		result.Type = models.Gauge
		currentValue, _ := strconv.ParseFloat(valueMEtric, 64)
		result.Value = currentValue
	}
	return result, true
}

func (p Postgres) GetList() map[string]models.Metric {
	// Запрос всех метрик из базы данных
	rows, err := p.Database.Query(`SELECT name, type, value FROM metrics`)
	if err != nil {
		return nil
	}
	defer rows.Close()

	result := make(map[string]models.Metric)

	for rows.Next() {
		var name, metricType, valueStr string

		// Считываем данные из строки
		if err := rows.Scan(&name, &metricType, &valueStr); err != nil {
			return nil
		}

		// Создаем метрику
		metric := models.Metric{
			Name: name,
			Type: models.MetricType(metricType),
		}

		// Преобразуем значение в зависимости от типа
		switch metric.Type {
		case models.Counter:
			value, err := strconv.ParseInt(valueStr, 10, 64)
			if err != nil {
				return nil
			}
			metric.Value = value

		case models.Gauge:
			value, err := strconv.ParseFloat(valueStr, 64)
			if err != nil {
				return nil
			}
			metric.Value = value

		default:
			return nil
		}

		// Добавляем метрику в результат
		result[name] = metric
	}

	// Проверяем ошибки, которые могли возникнуть при итерации
	if err := rows.Err(); err != nil {
		return nil
	}

	return result
}

func Init(connect string) error {
	var err error
	DB, err = sql.Open("pgx", connect)
	if err != nil {
		return err
	}
	DB.ExecContext(
		context.Background(),
		`CREATE TABLE IF NOT EXISTS metrics( 
    id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	type TEXT NOT NULL,
	value TEXT NOT NULL)`,
	)
	return err
}
