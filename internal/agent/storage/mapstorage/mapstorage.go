package mapstorage

import (
	"errors"

	"github.com/golovanevvs/metalecoll/internal/agent/model"
)

type Storage interface {
	SaveMetric(met model.Metric)
	GetMetricsMap() (map[string]model.Metric, error)
}
type aMemStorage struct {
	metrics map[string]model.Metric
}

func (ams *aMemStorage) SaveMetric(met model.Metric) {
	ams.metrics[met.Name] = met
}

func (ams *aMemStorage) GetMetricsMap() (map[string]model.Metric, error) {
	if len(ams.metrics) > 0 {
		return ams.metrics, nil
	}
	err := errors.New("нет сохранённых данных")
	return nil, err
}

func NewStorage() *aMemStorage {
	return &aMemStorage{
		metrics: make(map[string]model.Metric),
	}
}
