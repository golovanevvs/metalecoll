package server

import (
	"database/sql"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func PingDatabaseHandler(w http.ResponseWriter, r *http.Request, databaseDNS string) {
	srv.logger.Debugf("Открытие БД %v...", databaseDNS)
	db, err := sql.Open("pgx", databaseDNS)
	if err != nil {
		srv.logger.Errorf("Ошибка открытия БД: %v", err)
		srv.logger.Errorf("Отправлен код: %v", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()
	srv.logger.Debugf("Открытие БД прошло успешно")

	srv.logger.Debugf("Создание БД...")
	_, err = db.Exec("CREATE DATABASE metalecoll WITH OWNER = postgres ENCODING = 'UTF8' TABLESPACE = pg_default CONNECTION LIMIT = -1;")
	if err != nil {
		srv.logger.Errorf("Ошибка создания БД: %v", err)
		srv.logger.Errorf("Отправлен код: %v", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	srv.logger.Debugf("Создание БД прошло успешно")

	srv.logger.Debugf("Пингование БД...")
	if err := db.Ping(); err != nil {
		srv.logger.Errorf("Ошибка пингования БД: %v", err)
		srv.logger.Errorf("Отправлен код: %v", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	srv.logger.Debugf("Пингование БД прошло успешно")

	srv.logger.Debugf("Отправлен код: %v", http.StatusOK)
	w.WriteHeader(http.StatusOK)

}
