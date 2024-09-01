package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golovanevvs/metalecoll/internal/server/storage"
)

type metNameValue struct {
	BMetName  string `json:"name"`
	BMetValue string `json:"value"`
}

func GetMetricNamesHandler(w http.ResponseWriter, r *http.Request, store storage.Storage) {
	var metrics []metNameValue

	srv.logger.Debugf("")
	srv.logger.Debugf("Получение всех известных метрик из хранилища...")
	metricsMap := storage.GMs(store)
	srv.logger.Debugf("Получение всех известных метрик из хранилища прошло успешно")

	srv.logger.Debugf("")
	srv.logger.Debugf("Создание среза имя:значение...")
	for k, v := range metricsMap {
		value := fmt.Sprintf("%v", v.MetValue)
		m := metNameValue{
			BMetName:  k,
			BMetValue: value,
		}
		metrics = append(metrics, m)
	}
	srv.logger.Debugf("Создание среза имя:значение прошло успешно")
	srv.logger.Debugf("%v", metrics)

	srv.logger.Debugf("")
	srv.logger.Debugf("Сериализация данных в JSON...")
	resp, err := json.Marshal(metrics)
	if err != nil {
		srv.logger.Errorf("Ошибка сериализации данных в JSON")
		return
	}
	srv.logger.Debugf("Сериализация данных в JSON прошла успешно")

	srv.logger.Debugf("")
	srv.logger.Debugf("Устанавливаем заголовок Content-Type для передачи информации, кодированной в JSON")
	w.Header().Set("content-type", "application/json")

	srv.logger.Debugf("")
	srv.logger.Debugf("Отправляем код: %v", http.StatusOK)
	w.WriteHeader(http.StatusOK)

	srv.logger.Debugf("")
	srv.logger.Debugf("Отправка тела ответа...")
	w.Write(resp)
	srv.logger.Debugf("Отправка тела ответа прошла успешно")
}
