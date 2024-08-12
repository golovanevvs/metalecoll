package autil

import (
	"github.com/golovanevvs/metalecoll/internal/agent/storage/amapstorage"
	"github.com/golovanevvs/metalecoll/internal/server/model"
)

func SM(s amapstorage.AStorage, m model.Metric) {
	s.SaveMetric(m)
}

func GMM(s amapstorage.AStorage) (map[string]model.Metric, error) {
	if _, err := s.GetMetricsMap(); err != nil {
		return nil, err
	}
	return s.GetMetricsMap()
}
