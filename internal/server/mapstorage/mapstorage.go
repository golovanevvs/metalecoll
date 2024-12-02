// Модуль mapstorage предназначен для работы с map-хранилищем.
package mapstorage

import (
	"fmt"

	"github.com/golovanevvs/metalecoll/internal/server/model"
)

// Storage - интерфейс для работы с map-хранилищем.
type Storage interface {
	SaveMetricToMap(met model.Metric)
	GetMetricFromMap(name string) (model.Metric, error)
	GetMetricsFromMap() map[string]model.Metric
}

type memStorage struct {
	Metrics map[string]model.Metric
}

// NewMapStorage - конструктор map-хранилища.
func NewMapStorage() *memStorage {
	return &memStorage{
		Metrics: make(map[string]model.Metric),
	}
}

// SaveMetricToMap - сохранение метрики в хранилище.
func (ms *memStorage) SaveMetricToMap(met model.Metric) {
	ms.Metrics[met.MetName] = met
}

// GetMetricFromMap - получение метрики из хранилища.
func (ms *memStorage) GetMetricFromMap(name string) (model.Metric, error) {
	if _, inMap := ms.Metrics[name]; inMap {
		return ms.Metrics[name], nil
	}
	err := fmt.Errorf("в хранилище отсутствует запрошенный тип метрики")
	return model.Metric{}, err
}

// GetMetricsFromMap - получение всех метрик из хранилища.
func (ms *memStorage) GetMetricsFromMap() map[string]model.Metric {
	return ms.Metrics
}
