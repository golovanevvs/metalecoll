package sqlstorage

import (
	"context"
	"database/sql"

	"github.com/golovanevvs/metalecoll/internal/server/config"
	"github.com/golovanevvs/metalecoll/internal/server/storage/mapstorage"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type SQLStorage struct {
	name string
	db   *sql.DB
}

func New(c *config.Config) (*SQLStorage, error) {
	db, err := sql.Open("pgx", c.Storage.DatabaseDSN)
	if err != nil {
		return nil, err
	}
	return &SQLStorage{
		name: "БД Posgres " + c.Storage.DatabaseDSN,
		db:   db,
	}, nil
}

// GetNameDB возвращает название хранилища
func (s *SQLStorage) GetNameDB() string {
	return s.name
}

func (s *SQLStorage) SaveMetricsToDB(ctx context.Context, c *config.Config, mapStore mapstorage.Storage) error {
	return nil
}

func (s *SQLStorage) GetMetricsFromDB(ctx context.Context, c *config.Config) (mapstorage.Storage, error) {
	ms := mapstorage.New()
	return ms, nil
}

func (s *SQLStorage) Ping() error {
	if err := s.db.Ping(); err != nil {
		return err
	}
	return nil
}
