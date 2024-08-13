package mapstorage

import (
	"errors"

	"github.com/golovanevvs/metalecoll/internal/server/model"
)

type Storage interface {
	SaveMetric(met model.Metric)
	GetMetric(key string) (model.Metric, error)
	GetMetrics() map[string]model.Metric
}

type memStorage struct {
	Metrics map[string]model.Metric
}

func (ms *memStorage) SaveMetric(met model.Metric) {
	if ms.Metrics == nil {
		ms.Metrics = make(map[string]model.Metric)
	}
	ms.Metrics[met.MetType] = met
}

func (ms *memStorage) GetMetric(key string) (model.Metric, error) {
	if _, inMap := ms.Metrics[key]; inMap {
		return ms.Metrics[key], nil
	}
	err := errors.New("в хранилище отсутствует запрошенный тип метрики")
	return model.Metric{}, err
}

func (ms *memStorage) GetMetrics() map[string]model.Metric {
	return ms.Metrics
}

func NewStorage() *memStorage {
	return &memStorage{}
}
