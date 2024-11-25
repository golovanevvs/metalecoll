package service

import "github.com/golovanevvs/metalecoll/internal/server/model"

func (sv *getMetricsService) GetMetricsFromMap() map[string]model.Metric {
	return sv.mapStorage.GetMetricsFromMap()
}
