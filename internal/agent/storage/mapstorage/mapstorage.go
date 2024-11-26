package mapstorage

import (
	"errors"
	"sync"

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
	mu := new(sync.Mutex)
	mu.Lock()
	ams.metrics[met.Name] = met
	mu.Unlock()
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
