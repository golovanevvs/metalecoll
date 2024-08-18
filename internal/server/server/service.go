package server

import (
	"github.com/golovanevvs/metalecoll/internal/server/constants"
	"github.com/golovanevvs/metalecoll/internal/server/model"
	"github.com/golovanevvs/metalecoll/internal/server/util"
)

func procMetric(recMet model.Metric) *model.Metric {
	var newValue any
	switch recMet.MetType {
	case constants.GaugeType:
		newValue = recMet.MetValue.(float64)
	case constants.CounterType:
		if getValue, err := util.GM(srv.store, recMet.MetName); err != nil {
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
