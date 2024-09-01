package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golovanevvs/metalecoll/internal/server/constants"
	"github.com/golovanevvs/metalecoll/internal/server/storage"
)

func GetMetricValueHandler(w http.ResponseWriter, r *http.Request, store storage.Storage) {

	sbody := strings.Split(r.URL.Path, "/")

	if len(sbody) != 4 {
		srv.logger.Errorf("Структура тела запроса не соответствует ожидаемой. Получено тело запроса: %v", r.URL.Path)
		srv.logger.Errorf("")
		srv.logger.Errorf("Отправлен код: %v", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	srv.logger.Debugf("Чтение и разделение тела запроса прошло успешно")
	srv.logger.Debugf("")
	srv.logger.Debugf("Тело запроса: %v", r.URL.Path)

	srv.logger.Debugf("")
	srv.logger.Debugf("Параметры полученной метрики:")
	mM := sbody[1] // Тип метода
	srv.logger.Debugf("Тип метода: %v", mM)
	mT := sbody[2] // Тип метрики
	srv.logger.Debugf("Тип метрики: %v", mT)
	mN := sbody[3] // Имя метрики
	srv.logger.Debugf("Имя метрики: %v", mN)

	switch mT {
	case constants.GaugeType, constants.CounterType:
	default:
		srv.logger.Errorf("Неизвестный тип метрики")
		srv.logger.Errorf("")
		srv.logger.Errorf("Отправлен код: %v", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	srv.logger.Debugf("Проверка типа метрики прошла успешно")

	srv.logger.Debugf("")
	srv.logger.Debugf("Получение данных из хранилища...")
	metric, err := storage.GM(store, mN)
	if err != nil {
		fmt.Println(err)
		srv.logger.Errorf("Ошибка получения данных из хранилища")
		srv.logger.Errorf("Отправлен код: %v", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	srv.logger.Debugf("Получение данных из хранилища прошло успешно")

	value := fmt.Sprintf("%v", metric.MetValue)

	srv.logger.Debugf("")
	srv.logger.Debugf("Отправлен код: %v", http.StatusOK)
	w.WriteHeader(http.StatusOK)

	srv.logger.Debugf("")
	srv.logger.Debugf("Вывод полученных данных...")
	w.Write([]byte(value))
	srv.logger.Debugf("Вывод полученных данных прошёл успешно")
}
