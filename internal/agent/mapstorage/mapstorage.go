// Модуль mapstorage предназначен для работы с map-хранилищем.
package mapstorage

import (
	"errors"
	"sync"

	"github.com/golovanevvs/metalecoll/internal/agent/model"
)

// Storage интерфейс для работы с map-хранилищем.
type Storage interface {
	SaveMetric(met model.Metric)
	GetMetricsMap() (map[string]model.Metric, error)
}
type aMemStorage struct {
	metrics map[string]model.Metric
}

// SaveMetric сохраняет данные в map-хранилище.
func (ams *aMemStorage) SaveMetric(met model.Metric) {
	mu := new(sync.Mutex)
	mu.Lock()
	ams.metrics[met.Name] = met
	mu.Unlock()
}

// GetMetricsMap возвращает данные из map-хранилища.
func (ams *aMemStorage) GetMetricsMap() (map[string]model.Metric, error) {
	if len(ams.metrics) > 0 {
		return ams.metrics, nil
	}
	err := errors.New("нет сохранённых данных")
	return nil, err
}

// NewStorage - конструктор aMemStorage.
func NewStorage() *aMemStorage {
	return &aMemStorage{
		metrics: make(map[string]model.Metric),
	}
}
