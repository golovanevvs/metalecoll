package storage

import (
	"github.com/golovanevvs/metalecoll/internal/server/model"
)

type Storage interface {
	SaveMetric(met model.Metric)
	GetMetricsMap() (map[string]model.Metric, error)
}

func SM(s Storage, m model.Metric) {
	s.SaveMetric(m)
}

func GMM(s Storage) (map[string]model.Metric, error) {
	if _, err := s.GetMetricsMap(); err != nil {
		return nil, err
	}
	return s.GetMetricsMap()
}
