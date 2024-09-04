package server

import (
	"encoding/json"
	"net/http"

	"github.com/golovanevvs/metalecoll/internal/server/constants"
	"github.com/golovanevvs/metalecoll/internal/server/model"
	"github.com/golovanevvs/metalecoll/internal/server/service"
	"github.com/golovanevvs/metalecoll/internal/server/storage"
)

func UpdateMetricsJSONHandler(w http.ResponseWriter, r *http.Request, store storage.Storage) {
	var mVParse any
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
	// case constants.AContentTypeAJ:
	// default:
	// 	srv.logger.Errorf("Недопустимый content-type: %v", cT)
	// 	srv.logger.Errorf("Отправлен код: %v", http.StatusBadRequest)
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	// srv.logger.Debugf("Проверка Content-Type прошла успешно")

	srv.logger.Debugf("Проверка типа метода...")

	srv.logger.Debugf("Проверка наличия имени метрики...")

	if req.ID == "" {
		srv.logger.Errorf("Имя метрики не задано")
		srv.logger.Errorf("Отправлен код: %v", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	srv.logger.Debugf("Проверка наличия имени метрики прошла успешно")

	switch req.MType {
	case constants.GaugeType:
		mVParse = *req.Value
	case constants.CounterType:
		mVParse = *req.Delta
	default:
		srv.logger.Errorf("Неизвестный тип метрики")
		srv.logger.Errorf("Отправлен код: %v", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	receivedMetric := model.Metric{MetType: req.MType, MetName: req.ID, MetValue: mVParse}

	srv.logger.Debugf("Обновление метрики...")
	calcMetric := service.ProcMetric(receivedMetric, store)
	srv.logger.Debugf("%v", calcMetric)
	srv.logger.Debugf("Обновление метрики прошло успешно")

	srv.logger.Debugf("Обновление хранилища...")
	storage.SM(store, *calcMetric)
	srv.logger.Debugf("Обновление хранилища прошло успешно")

	srv.logger.Debugf("Обновлённое хранилище: %v", storage.GMs(store))

	srv.logger.Debugf("Формирование тела ответа...")
	switch req.MType {
	case constants.GaugeType:
		if v, ok := calcMetric.MetValue.(float64); ok {
			resp = Metrics{
				ID:    calcMetric.MetName,
				MType: calcMetric.MetType,
				Value: &v,
			}
		}
	case constants.CounterType:
		if v, ok := calcMetric.MetValue.(int64); ok {
			resp = Metrics{
				ID:    calcMetric.MetName,
				MType: calcMetric.MetType,
				Delta: &v,
			}
		}
	}
	srv.logger.Debugf("Формирование тела ответа прошло успешно")

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
		return
	}
	srv.logger.Debugf("Кодирование в JSON прошло успешно")

	srv.logger.Debugf("Отправлен Content-Type: %v", constants.ContentTypeAJ)
	w.Header().Set("Content-Type", "application/json")

	srv.logger.Debugf("Отправлен код: %v", http.StatusOK)
	w.WriteHeader(http.StatusOK)
	w.Write(enc)
}
