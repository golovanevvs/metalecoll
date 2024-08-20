package util

import (
	"fmt"

	"github.com/golovanevvs/metalecoll/internal/server/model"
	"github.com/golovanevvs/metalecoll/internal/server/storage/mapstorage"
)

func SM(s mapstorage.Storage, m model.Metric) {
	fmt.Println("Запуск SM")
	s.SaveMetric(m)
	fmt.Println("Завершение SM")
}

func GM(s mapstorage.Storage, key string) (model.Metric, error) {
	if _, err := s.GetMetric(key); err != nil {
		return model.Metric{}, err
	}
	return s.GetMetric(key)
}

func GMs(s mapstorage.Storage) map[string]model.Metric {
	return s.GetMetrics()
}
