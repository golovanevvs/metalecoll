package service

import (
	"fmt"

	"github.com/golovanevvs/metalecoll/internal/server/constants"
	"github.com/golovanevvs/metalecoll/internal/server/model"
)

func (sv *updateMetricsService) UpdateMetric(recMet model.Metric) *model.Metric {
	updatedMetric := sv.procMetric(recMet)
	sv.mapStorage.SaveMetricToMap(*updatedMetric)
	return updatedMetric
}

func (sv *updateMetricsService) procMetric(recMet model.Metric) *model.Metric {
	var newValue any

	switch recMet.MetType {

	case constants.GaugeType:
		newValue = recMet.MetValue.(float64)

	case constants.CounterType:
		if getValue, err := sv.mapStorage.GetMetricFromMap(recMet.MetName); err != nil {
			newValue = recMet.MetValue.(int64)
		} else {
			newValue = (getValue.MetValue.(int64)) + recMet.MetValue.(int64)
		}
	}

	slice := make([]int, 0, 10000000)
	for i := 0; i < 10000000; i++ {
		slice = append(slice, i)
	}
	fmt.Println(len(slice))

	return &model.Metric{
		MetType:  recMet.MetType,
		MetName:  recMet.MetName,
		MetValue: newValue,
	}
}
