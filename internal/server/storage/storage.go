package storage

import (
	"github.com/golovanevvs/metalecoll/internal/server/model"
)

type Storage interface {
	SaveMetric(met model.Metric)
	GetMetric(key string) (model.Metric, error)
	GetMetrics() map[string]model.Metric
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
