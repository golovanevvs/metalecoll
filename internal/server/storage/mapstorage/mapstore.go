package mapstorage

import "github.com/golovanevvs/metalecoll/internal/server/model"

// type Repositories interface {
// 	metricSave(m *model.Metric)
// }

type MemStorage struct {
	Metrics map[string]model.Metric
}

// func updateStorage(s *MemStorage, m model.Metric) {
// 	s.metricSave(m)
// }
