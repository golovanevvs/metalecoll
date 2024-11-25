package postgres

import (
	"context"

	"github.com/golovanevvs/metalecoll/internal/server/config"
	"github.com/golovanevvs/metalecoll/internal/server/constants"
	"github.com/golovanevvs/metalecoll/internal/server/mapstorage"
	"github.com/golovanevvs/metalecoll/internal/server/model"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type allPostgres struct {
	name string
	db   *sqlx.DB
}

// NewPostgres - конструктор
func NewPostgres(databaseDSN string) (*allPostgres, error) {
	// открытие БД
	db, err := sqlx.Open("pgx", databaseDSN)
	if err != nil {
		return nil, err
	}

	// пингование БД
	if err = db.Ping(); err != nil {
		return nil, err
	}

	// создание экземпляра allPostgres
	st := &allPostgres{
		name: "БД Postgres " + databaseDSN,
		db:   db,
	}

	return st, nil
}

// GetNameDB возвращает название хранилища
func (s *allPostgres) GetNameDB() string {
	return s.name
}

// SaveMetricsToDB сохраняет метрики в БД
func (s *allPostgres) SaveMetricsToDB(ctx context.Context, c *config.Config, mapStore mapstorage.Storage) error {
	s.db.ExecContext(ctx, `
	TRUNCATE TABLE metrics RESTART IDENTITY;
	`)
	mapMetrics := mapStore.GetMetricsFromMap()
	for _, metric := range mapMetrics {
		switch metric.MetType {
		case constants.GaugeType:
			_, err := s.db.ExecContext(ctx, `
			INSERT INTO metrics
				(metric_name, metric_type, gauge_value)
				VALUES
				($1, $2, $3);
			`, metric.MetName, metric.MetType, metric.MetValue)
			if err != nil {
				return err
			}
		case constants.CounterType:
			_, err := s.db.ExecContext(ctx, `
			INSERT INTO metrics
				(metric_name, metric_type, counter_value)
				VALUES
				($1, $2, $3);
			`, metric.MetName, metric.MetType, metric.MetValue)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// GetMetricsFromDB получает метрики из БД
func (s *allPostgres) GetMetricsFromDB(ctx context.Context, c *config.Config) (mapstorage.Storage, error) {
	var (
		gaugeValue   float64
		counterValue int64
	)
	ms := mapstorage.NewMapStorage()
	rows, err := s.db.QueryContext(ctx, `
	SELECT
		metric_name,
		metric_type,
		gauge_value,
		counter_value
	FROM metrics
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m model.Metric
		if err := rows.Scan(&m.MetName, &m.MetType, &gaugeValue, &counterValue); err != nil {
			return nil, err
		}
		switch m.MetType {
		case constants.GaugeType:
			m.MetValue = gaugeValue
		case constants.CounterType:
			m.MetValue = counterValue
		}
		ms.SaveMetricToMap(m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ms, nil
}

func (s *allPostgres) Ping() error {
	if err := s.db.Ping(); err != nil {
		return err
	}
	return nil
}

func (s *allPostgres) CloseDB() error {
	return s.db.Close()
}
