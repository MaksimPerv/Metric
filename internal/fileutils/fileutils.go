package fileutils

import (
	"encoding/json"
	"github.com/MaksimPerv/Metric/internal/models"
	"os"
)

type Producer struct {
	file   *os.File
	encode *json.Encoder
}

var File *Producer

func NewProducer(name string, flag bool) (*Producer, error) {
	var pr int
	if flag {
		pr = os.O_WRONLY | os.O_CREATE | os.O_APPEND | os.O_TRUNC
	} else {
		pr = os.O_WRONLY | os.O_CREATE | os.O_APPEND
	}
	file, err := os.OpenFile(name, pr, 0666)
	if err != nil {
		return nil, err
	}
	return &Producer{
		file:   file,
		encode: json.NewEncoder(file),
	}, nil
}

func (producer *Producer) Write(data map[string]models.Metric) error {
	return producer.encode.Encode(data)
}
func (producer *Producer) WriteOne(data models.Metric) error {
	return producer.encode.Encode(data)
}
