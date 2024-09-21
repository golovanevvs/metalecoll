package storage

import (
	"context"

	"github.com/golovanevvs/metalecoll/internal/server/config"
	"github.com/golovanevvs/metalecoll/internal/server/storage/filestorage"
	"github.com/golovanevvs/metalecoll/internal/server/storage/mapstorage"
	"github.com/golovanevvs/metalecoll/internal/server/storage/sqlstorage"
)

// StorageDB - интерфейс работы с основным хранилищем
type StorageDB interface {
	GetNameDB() string
	SaveMetricsToDB(ctx context.Context, c *config.Config, mapStore mapstorage.Storage) error
	GetMetricsFromDB(ctx context.Context, c *config.Config) (mapstorage.Storage, error)
	Ping() error
}

// New - выбор основного хранилища: если флаг (-d) пуст, то выбирается файловое хранилище, иначе - БД
func New(c *config.Config) (StorageDB, error) {
	switch c.Storage.DatabaseDSN {
	case "":
		s := filestorage.New(c)
		return s, nil
	default:
		s, err := sqlstorage.New(c)
		if err != nil {
			return nil, err
		}
		return s, nil
	}
}
