package handler

// import (
// 	"net/http"

// 	_ "github.com/jackc/pgx/v5/stdlib"
// )

// func (s server) PingDatabaseHandler(w http.ResponseWriter, r *http.Request) {
// 	//ctx := r.Context()

// 	srv.logger.Debugf("Открытие БД прошло успешно")

// 	srv.logger.Debugf("Пингование БД...")
// 	if err := s.dbStore.Ping(); err != nil {
// 		srv.logger.Errorf("Ошибка пингования БД: %v", err)
// 		srv.logger.Errorf("Отправлен код: %v", http.StatusInternalServerError)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	srv.logger.Debugf("Пингование БД прошло успешно")

// 	srv.logger.Debugf("Отправлен код: %v", http.StatusOK)
// 	w.WriteHeader(http.StatusOK)

// }
