package fileutils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/MaksimPerv/Metric/internal/models"
	"io"
	"os"
)

type Producer struct {
	file   *os.File
	encode *json.Encoder
}

var File *Producer

func NewProducer(name string) (*Producer, error) {

	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &Producer{
		file:   file,
		encode: json.NewEncoder(file),
	}, nil
}

func (producer *Producer) Write(data map[string]models.Metric) error {
	for _, value := range data {
		err := producer.encode.Encode(value)
		if err != nil {
			return err
		}
	}
	return nil
}
func (producer *Producer) Close() {
	producer.file.Close()
	return
}

func (producer *Producer) Read() map[string]models.Metric {
	result := make(map[string]models.Metric)
	// Сохраняем текущую позицию записи
	currentPos, err := producer.file.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil
	}

	// Перемещаемся в начало для чтения
	if _, err := producer.file.Seek(0, io.SeekStart); err != nil {
		return nil
	}

	// Восстанавливаем позицию после чтения
	defer func() {
		_, _ = producer.file.Seek(currentPos, io.SeekStart)
	}()
	scanner := bufio.NewScanner(producer.file)
	for scanner.Scan() {
		line := scanner.Text()
		var metric models.Metric
		if err := json.Unmarshal([]byte(line), &metric); err != nil {
			fmt.Printf("Ошибка декодирования: %v\n", err)
			continue
		}
		result[metric.Name] = metric
	}
	return result
}

func (producer *Producer) WriteOne(data models.Metric) error {
	return producer.encode.Encode(data)
}
