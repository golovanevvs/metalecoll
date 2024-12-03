// Модуль service содержит в себе бизнес-логику приложения.
package service

import (
	"github.com/golovanevvs/metalecoll/internal/server/mapstorage"
	"github.com/golovanevvs/metalecoll/internal/server/model"
	"github.com/golovanevvs/metalecoll/internal/server/storage"
)

// IUpdateMetricsService - интерфейс сервиса обновления метрик.
type IUpdateMetricsService interface {
	UpdateMetric(recMet model.Metric) *model.Metric
}

// IGetMetricsService - интерфейс сервиса получения метрик.
type IGetMetricsService interface {
	GetMetricFromMap(name string) (model.Metric, error)
	GetMetricsFromMap() map[string]model.Metric
}

// IPingDatabaseService - интерфейс сервиса проверки соединения с базой данных.
type IPingDatabaseService interface {
	Ping() error
}

type updateMetricsService struct {
	mapStorage mapstorage.Storage
	st         storage.IStorageDB
}

type getMetricsService struct {
	mapStorage mapstorage.Storage
	st         storage.IStorageDB
}

type pingDatabaseService struct {
	st storage.IStorageDB
}

type Service struct {
	IUpdateMetricsService
	IGetMetricsService
	IPingDatabaseService
}

// NewUpdateMetricsService - конструктор сервиса обновления метрик.
func NewUpdateMetricsService(mst mapstorage.Storage, st storage.IStorageDB) *updateMetricsService {
	return &updateMetricsService{
		mapStorage: mst,
		st:         st,
	}
}

// NewGetMetricsService - конструктор сервиса получения метрик.
func NewGetMetricsService(mst mapstorage.Storage, st storage.IStorageDB) *getMetricsService {
	return &getMetricsService{
		mapStorage: mst,
		st:         st,
	}
}

// NewPingDatabaseService - конструктор сервиса проверки соединения с базой данных.
func NewPingDatabaseService(st storage.IStorageDB) *pingDatabaseService {
	return &pingDatabaseService{st: st}
}

// NewService - конструктор сервиса.
func NewService(mst mapstorage.Storage, st storage.IStorageDB) *Service {
	return &Service{
		IUpdateMetricsService: NewUpdateMetricsService(mst, st),
		IGetMetricsService:    NewGetMetricsService(mst, st),
		IPingDatabaseService:  NewPingDatabaseService(st),
	}
}
