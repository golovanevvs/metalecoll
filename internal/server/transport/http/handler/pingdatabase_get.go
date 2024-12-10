package handler

import (
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// PingDatabase - Пингование БД.
func (hd *handler) PingDatabase(w http.ResponseWriter, r *http.Request) {
	hd.lg.Debugf("Открытие БД прошло успешно")

	hd.lg.Debugf("Пингование БД...")
	if err := hd.sv.Ping(); err != nil {
		hd.lg.Errorf("Ошибка пингования БД: %v", err)
		hd.lg.Errorf("Отправлен код: %v", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	hd.lg.Debugf("Пингование БД прошло успешно")

	hd.lg.Debugf("Отправлен код: %v", http.StatusOK)
	w.WriteHeader(http.StatusOK)

}
