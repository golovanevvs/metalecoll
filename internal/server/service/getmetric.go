package service

import "github.com/golovanevvs/metalecoll/internal/server/model"

func (sv *getMetricsService) GetMetricFromMap(name string) (model.Metric, error) {
	return sv.mapStorage.GetMetricFromMap(name)
}
