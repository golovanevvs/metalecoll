package service

import (
	"github.com/golovanevvs/metalecoll/internal/server/constants"
	"github.com/golovanevvs/metalecoll/internal/server/model"
	"github.com/golovanevvs/metalecoll/internal/server/storage"
)

func ProcMetric(recMet model.Metric, store storage.Storage) *model.Metric {
	var newValue any

	switch recMet.MetType {

	case constants.GaugeType:
		newValue = recMet.MetValue.(float64)

	case constants.CounterType:
		if getValue, err := store.GetMetric(recMet.MetName); err != nil {
			newValue = recMet.MetValue.(int64)
		} else {
			newValue = getValue.MetValue.(int64) + recMet.MetValue.(int64)
		}
	}

	return &model.Metric{
		MetType:  recMet.MetType,
		MetName:  recMet.MetName,
		MetValue: newValue,
	}
}
