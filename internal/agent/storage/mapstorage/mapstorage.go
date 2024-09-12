package mapstorage

import (
	"errors"

	"github.com/golovanevvs/metalecoll/internal/server/model"
)

type aMemStorage struct {
	metrics map[string]model.Metric
}

func (ams *aMemStorage) SaveMetric(met model.Metric) {
	if ams.metrics == nil {
		ams.metrics = make(map[string]model.Metric)
	}
	ams.metrics[met.MetName] = met
}

func (ams *aMemStorage) GetMetricsMap() (map[string]model.Metric, error) {
	if len(ams.metrics) > 0 {
		return ams.metrics, nil
	}
	err := errors.New("нет сохранённых данных")
	return nil, err
}

func NewStorage() *aMemStorage {
	return &aMemStorage{}
}
