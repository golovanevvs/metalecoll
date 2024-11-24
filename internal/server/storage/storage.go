package storage

import (
	"context"

	"github.com/golovanevvs/metalecoll/internal/server/config"
	"github.com/golovanevvs/metalecoll/internal/server/storage/filestorage"
	"github.com/golovanevvs/metalecoll/internal/server/storage/mapstorage"
	"github.com/golovanevvs/metalecoll/internal/server/storage/postgres"
)

// StorageDB - интерфейс работы с основным хранилищем
type StorageDB interface {
	GetNameDB() string
	SaveMetricsToDB(ctx context.Context, c *config.Config, mapStore mapstorage.Storage) error
	GetMetricsFromDB(ctx context.Context, c *config.Config) (mapstorage.Storage, error)
	Ping() error
}

// NewStorage - выбор основного хранилища: если флаг (-d) пуст, то выбирается файловое хранилище, иначе - БД
func NewStorage(c *config.Config) (StorageDB, error) {
	switch c.Storage.DatabaseDSN {
	case "":
		s := filestorage.NewFileStorage(c.Storage.FileStoragePath)
		return s, nil
	default:
		s, err := postgres.NewPostgres(c.Storage.DatabaseDSN)
		if err != nil {
			return nil, err
		}
		return s, nil
	}
}
