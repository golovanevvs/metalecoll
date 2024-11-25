package storage

import (
	"context"

	"github.com/golovanevvs/metalecoll/internal/server/config"
	"github.com/golovanevvs/metalecoll/internal/server/mapstorage"
)

// IStorageDB - интерфейс работы с основным хранилищем
type IStorageDB interface {
	GetNameDB() string
	SaveMetricsToDB(ctx context.Context, c *config.Config, mapStore mapstorage.Storage) error
	GetMetricsFromDB(ctx context.Context, c *config.Config) (mapstorage.Storage, error)
	Ping() error
	CloseDB() error
}

type StorageDB struct {
	IStorageDB
}

func NewStorage(mainStore IStorageDB) *StorageDB {
	return &StorageDB{
		IStorageDB: mainStore,
	}
}
