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
	err := errors.New("В хранилище отсутствует запрошенный тип метрики")
	return model.Metric{}, err
}

func (ms *memStorage) GetMetrics() map[string]model.Metric {
	return ms.Metrics
}

func (ms *memStorage) NewStorage() Storage {
	return ms
}

func NewStorage() *memStorage {
	return &memStorage{}
}

func SM(s Storage, m model.Metric) {
	s.SaveMetric(m)
}

func GM(s Storage, key string) (model.Metric, error) {
	if _, err := s.GetMetric(key); err != nil {
		return model.Metric{}, err
	}
	return s.GetMetric(key)
}

func GMs(s Storage) map[string]model.Metric {
	return s.GetMetrics()
}
