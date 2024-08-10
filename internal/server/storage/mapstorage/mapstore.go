package mapstorage

import "github.com/golovanevvs/metalecoll/internal/server/model"

type MemStorage struct {
	Metrics map[string]model.Metric
}

func (ms *MemStorage) SaveMetric(met model.Metric) {
	if ms.Metrics == nil {
		ms.Metrics = make(map[string]model.Metric)
	}
	ms.Metrics[met.MetType] = met
}
