package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/golovanevvs/metalecoll/internal/server/config"
	"github.com/golovanevvs/metalecoll/internal/server/model"
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

	st := &SQLStorage{
		name: "БД Posgres " + c.Storage.DatabaseDSN,
		db:   db,
	}

	tableMetrics := "metrics"

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	exist, err := st.tablesExist(tableMetrics)
	if err != nil {
		return nil, err
	}

	if !exist {
		_, err = db.ExecContext(ctx, `
		CREATE TABLE metrics (
			id SERIAL PRIMARY KEY,
			metric_name VARCHAR(250) NOT NULL,
			metric_type VARCHAR(250) NOT NULL,
			gauge_value DOUBLE PRECISION DEFAULT NULL,
			counter_value INTEGER DEFAULT NULL
		)
		`,
		)
		if err != nil {
			return nil, err
		}
		exist2, err := st.tablesExist(tableMetrics)
		if err != nil {
			return nil, err
		}
		if !exist2 {
			return nil, errors.New("неизвестная ошибка создания таблицы metrics")
		}
	}
	return st, nil
}

func (s SQLStorage) tablesExist(nameTable string) (bool, error) {
	var exists bool
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	row := s.db.QueryRowContext(ctx, `
		SELECT EXISTS
			(SELECT FROM information_schema.tables
			WHERE table_name = 'metrics')
		`,
	)

	err := row.Scan(&exists)
	if err != nil {
		return false, err
	}

	if exists {
		return true, nil
	}

	return false, nil
}

// GetNameDB возвращает название хранилища
func (s *SQLStorage) GetNameDB() string {
	return s.name
}

func (s *SQLStorage) SaveMetricsToDB(ctx context.Context, c *config.Config, mapStore mapstorage.Storage) error {
	s.db.ExecContext(ctx, `
	TRUNCATE TABLE metrics RESTART IDENTITY;
	`)
	mapMetrics := mapStore.GetMetrics()
	for _, metric := range mapMetrics {
		switch metric.MetType {
		case c.MetricTypeNames.GaugeType:
			_, err := s.db.ExecContext(ctx, `
			INSERT INTO metrics
				(metric_name, metric_type, gauge_value)
				VALUES
				($1, $2, $3);
			`, metric.MetName, metric.MetType, metric.MetValue)
			if err != nil {
				return err
			}
		case c.MetricTypeNames.CounterType:
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

func (s *SQLStorage) GetMetricsFromDB(ctx context.Context, c *config.Config) (mapstorage.Storage, error) {
	var (
		gaugeValue   float64
		counterValue int64
	)
	ms := mapstorage.New()
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
		case c.MetricTypeNames.GaugeType:
			m.MetValue = gaugeValue
		case c.MetricTypeNames.CounterType:
			m.MetValue = counterValue
		}
		ms.SaveMetric(m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ms, nil
}

func (s *SQLStorage) Ping() error {
	if err := s.db.Ping(); err != nil {
		return err
	}
	return nil
}
