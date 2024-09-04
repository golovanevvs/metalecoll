package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golovanevvs/metalecoll/internal/server/constants"
	"github.com/golovanevvs/metalecoll/internal/server/storage"
)

func GetMetricValueJSONHandler(w http.ResponseWriter, r *http.Request, store storage.Storage) {
	var req, resp Metrics

	srv.logger.Debugf("Декодирование JSON...")

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		srv.logger.Errorf("Ошибка декодирования JSON: %v", err)
		srv.logger.Errorf("Отправлен код: %v", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	srv.logger.Debugf("Декодирование JSON прошло успешно")

	// srv.logger.Debugf("Проверка Content-Type...")
	// cT := r.Header.Get("Content-Type")
	// switch cT {
	// case constants.ContentTypeAJ:
	// default:
	// 	srv.logger.Errorf("Недопустимый content-type: %v", cT)
	// 	srv.logger.Errorf("Отправлен код: %v", http.StatusBadRequest)
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	// srv.logger.Debugf("Проверка Content-Type прошла успешно")

	srv.logger.Debugf("Получение данных из хранилища по name %v...", req.ID)
	metric, err := storage.GM(store, req.ID)
	if err != nil {
		fmt.Println(err)
		srv.logger.Errorf("Ошибка получения данных из хранилища")
		srv.logger.Errorf("Отправлен код: %v", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	srv.logger.Debugf("Получение данных из хранилища прошло успешно")

	srv.logger.Debugf("Формирование тела ответа...")
	switch req.MType {
	case constants.GaugeType:
		if v, ok := metric.MetValue.(float64); ok {
			resp = Metrics{
				ID:    metric.MetName,
				MType: metric.MetType,
				Value: &v,
			}
		}
	case constants.CounterType:
		if v, ok := metric.MetValue.(int64); ok {
			resp = Metrics{
				ID:    metric.MetName,
				MType: metric.MetType,
				Delta: &v,
			}
		}
	}
	srv.logger.Debugf("Формирование тела ответа прошло успешно")

	srv.logger.Debugf("Отправлен Content-Type: %v", constants.ContentTypeAJ)
	w.Header().Set("Content-Type", constants.ContentTypeAJ)

	srv.logger.Debugf("Кодирование в JSON...")
	// enc := json.NewEncoder(w)
	// if err := enc.Encode(resp); err != nil {
	// 	srv.logger.Errorf("Ошибка кодирования JSON: %v", err)
	// 	srv.logger.Errorf("Отправлен код: %v", http.StatusInternalServerError)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	enc, err := json.Marshal(resp)
	if err != nil {
		srv.logger.Errorf("Ошибка кодирования: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	srv.logger.Debugf("Кодирование в JSON прошло успешно")

	srv.logger.Debugf("Отправлен код: %v", http.StatusOK)
	w.WriteHeader(http.StatusOK)
	w.Write(enc)
}
