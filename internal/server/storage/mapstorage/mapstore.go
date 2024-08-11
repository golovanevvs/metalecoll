package mapstorage

import (
	"errors"

	"github.com/golovanevvs/metalecoll/internal/server/model"
)

type Storage interface {
	SaveMetric(met model.Metric)
	GetMetric(key string) model.Metric
	NewStorage() Storage
}

type MemStorage struct {
	Metrics map[string]model.Metric
}

func (ms *MemStorage) SaveMetric(met model.Metric) {
	if ms.Metrics == nil {
		ms.Metrics = make(map[string]model.Metric)
	}
	ms.Metrics[met.MetType] = met
}

func (ms *MemStorage) GetMetric(key string) model.Metric {
	return ms.Metrics[key]
}

func (ms *MemStorage) NewStorage() Storage {
	return ms
}

func NewStorage() (a Storage) {
	return a.NewStorage()
}

func SM(s Storage, m model.Metric) {
	s.SaveMetric(m)
}

func GM(s Storage, key string) (model.Metric, error) {
	if _, inMap := s[key]; inMap {
		return s.GetMetric(key), nil
	}
	err := errors.New("Данные в хранилище отсутствуют")
	return nil, err
}
