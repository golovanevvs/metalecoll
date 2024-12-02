package service

import "github.com/golovanevvs/metalecoll/internal/server/model"

// GetMetricsFromMap получение всех метрик из map-хранилища.
func (sv *getMetricsService) GetMetricsFromMap() map[string]model.Metric {
	return sv.mapStorage.GetMetricsFromMap()
}
