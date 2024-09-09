package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/golovanevvs/metalecoll/internal/server/storage"
)

func GetMetricValueHandler(w http.ResponseWriter, r *http.Request, store storage.Storage) {
	srv.logger.Debugf("Тело запроса: %v", r.URL.Path)

	srv.logger.Debugf("Параметры полученной метрики:")
	mT := chi.URLParam(r, "type")
	srv.logger.Debugf("Тип метрики: %v", mT)
	mN := chi.URLParam(r, "name")
	srv.logger.Debugf("Имя метрики: %v", mN)
	mV := chi.URLParam(r, "value")
	srv.logger.Debugf("Значение метрики: %v", mV)

	srv.logger.Debugf("Получение данных из хранилища...")
	metric, err := store.GetMetric(mN)
	if err != nil {
		fmt.Println(err)
		srv.logger.Errorf("Ошибка получения данных из хранилища")
		srv.logger.Errorf("Отправлен код: %v", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	srv.logger.Debugf("Получение данных из хранилища прошло успешно")

	value := fmt.Sprintf("%v", metric.MetValue)

	srv.logger.Debugf("Отправлен код: %v", http.StatusOK)
	w.WriteHeader(http.StatusOK)

	srv.logger.Debugf("Вывод полученных данных...")
	w.Write([]byte(value))
	srv.logger.Debugf("Вывод полученных данных прошёл успешно")
}
