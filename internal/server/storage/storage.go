package storage

import (
	"github.com/golovanevvs/metalecoll/internal/server/model"
)

type Storage interface {
	SaveMetric(met model.Metric)
	GetMetric(key string) (model.Metric, error)
	GetMetrics() map[string]model.Metric
}

type StorageDB interface {
	SaveToDB(m *model.Metric) error
	GetFromDB(name string) (model.Metric, error)
	Ping() error
}
