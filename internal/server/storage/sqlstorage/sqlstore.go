package sqlstorage

import "database/sql"

type SQLDB struct {
	db *sql.DB
}

func New(db *sql.DB) *SQLDB {
	return &SQLDB{
		db: db,
	}
}
