package sqlstorage

import (
	"database/sql"

	"github.com/golovanevvs/metalecoll/internal/server/model"
)

type SQLDB struct {
	db *sql.DB
}

func New(db *sql.DB) *SQLDB {
	return &SQLDB{
		db: db,
	}
}

func (s *SQLDB) SaveToDB(m *model.Metric) error {
	return nil
}

func (s *SQLDB) GetFromDB(name string) (model.Metric, error) {
	return model.Metric{}, nil
}

func (s *SQLDB) Ping() error {
	if err := s.db.Ping(); err != nil {
		return err
	}
	return nil
}
