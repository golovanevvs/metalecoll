package mapstorage

import (
	"github.com/golovanevvs/metalecoll/internal/server/model"
	"github.com/golovanevvs/metalecoll/internal/server/server"
)

type MetricRepository struct {
	metric *model.Metric
}

// type ii interface {
// 	MetricProc(s MemStorage)
// }

func (s *MemStorage) MetricSave(m *MetricRepository) {
	if s.Metrics == nil {
		s.Metrics = make(map[string]model.Metric)
	}
	s.Metrics[m.metric.MetType] = *m.metric
}

func (m *MetricRepository) MetricProc(s *MemStorage) *model.Metric {
	var newValue any
	switch m.metric.MetType {
	case server.GaugeType:
		newValue = m.metric.MetValue.(float64)
	case server.CounterType:
		if value, inMap := s.Metrics[server.CounterType]; inMap {
			newValue = value.MetValue.(int64) + m.metric.MetValue.(int64)
		} else {
			newValue = m.metric.MetValue.(int64)
		}
	}
	return &model.Metric{m.metric.MetType, m.metric.MetName, newValue}
}
