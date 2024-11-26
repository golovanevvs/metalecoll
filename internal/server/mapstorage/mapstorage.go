package mapstorage

import (
	"errors"

	"github.com/golovanevvs/metalecoll/internal/server/model"
)

type Storage interface {
	SaveMetricToMap(met model.Metric)
	GetMetricFromMap(name string) (model.Metric, error)
	GetMetricsFromMap() map[string]model.Metric
}

type memStorage struct {
	Metrics map[string]model.Metric
}

func NewMapStorage() *memStorage {
	return &memStorage{}
}

func (ms *memStorage) SaveMetricToMap(met model.Metric) {
	if ms.Metrics == nil {
		ms.Metrics = make(map[string]model.Metric)
	}
	ms.Metrics[met.MetName] = met
}

func (ms *memStorage) GetMetricFromMap(name string) (model.Metric, error) {
	if _, inMap := ms.Metrics[name]; inMap {
		return ms.Metrics[name], nil
	}
	err := errors.New("в хранилище отсутствует запрошенный тип метрики")
	return model.Metric{}, err
}

func (ms *memStorage) GetMetricsFromMap() map[string]model.Metric {
	return ms.Metrics
}
