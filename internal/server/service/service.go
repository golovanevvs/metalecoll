package service

import (
	"github.com/golovanevvs/metalecoll/internal/server/mapstorage"
	"github.com/golovanevvs/metalecoll/internal/server/model"
	"github.com/golovanevvs/metalecoll/internal/server/storage"
)

type IUpdateMetricsService interface {
	UpdateMetric(recMet model.Metric) *model.Metric
}

type IGetMetricsService interface {
	GetMetricFromMap(name string) (model.Metric, error)
	GetMetricsFromMap() map[string]model.Metric
}

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

type PingDatabaseService struct {
	st storage.IStorageDB
}

type Service struct {
	IUpdateMetricsService
	IGetMetricsService
	IPingDatabaseService
}

func NewUpdateMetricsService(mst mapstorage.Storage, st storage.IStorageDB) *updateMetricsService {
	return &updateMetricsService{
		mapStorage: mst,
		st:         st,
	}
}

func NewGetMetricsService(mst mapstorage.Storage, st storage.IStorageDB) *getMetricsService {
	return &getMetricsService{
		mapStorage: mst,
		st:         st,
	}
}

func NewPingDatabaseService(st storage.IStorageDB) *PingDatabaseService {
	return &PingDatabaseService{st: st}
}

func NewService(mst mapstorage.Storage, st storage.IStorageDB) *Service {
	return &Service{
		IUpdateMetricsService: NewUpdateMetricsService(mst, st),
		IGetMetricsService:    NewGetMetricsService(mst, st),
		IPingDatabaseService:  NewPingDatabaseService(st),
	}
}
