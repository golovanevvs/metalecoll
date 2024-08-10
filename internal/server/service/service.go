package service

import (
	"github.com/golovanevvs/metalecoll/internal/server/constants"
	"github.com/golovanevvs/metalecoll/internal/server/model"
	"github.com/golovanevvs/metalecoll/internal/server/storage/mapstorage"
)

func ProcMetric(met model.Metric, ms mapstorage.MemStorage) *model.Metric {
	var newValue any
	switch met.MetType {
	case constants.GaugeType:
		newValue = met.MetValue.(float64)
	case constants.CounterType:
		if value, inMap := ms.Metrics[constants.CounterType]; inMap {
			newValue = value.MetValue.(int64) + met.MetValue.(int64)
		} else {
			newValue = met.MetValue.(int64)
		}
	}
	return &model.Metric{
		MetType:  met.MetType,
		MetName:  met.MetName,
		MetValue: newValue,
	}
}
