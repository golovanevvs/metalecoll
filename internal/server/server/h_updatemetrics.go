package server

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/golovanevvs/metalecoll/internal/server/constants"
	"github.com/golovanevvs/metalecoll/internal/server/model"
	"github.com/golovanevvs/metalecoll/internal/server/service"
	"github.com/golovanevvs/metalecoll/internal/server/storage"
)

func UpdateMetricsHandler(w http.ResponseWriter, r *http.Request, store storage.Storage) {
	var mVParse any
	var err error

	srv.logger.Debugf("Тело запроса: %v", r.URL.Path)

	srv.logger.Debugf("Параметры полученной метрики:")
	mT := chi.URLParam(r, "type")
	srv.logger.Debugf("Тип метрики: %v", mT)
	mN := chi.URLParam(r, "name")
	srv.logger.Debugf("Имя метрики: %v", mN)
	mV := chi.URLParam(r, "value")
	srv.logger.Debugf("Значение метрики: %v", mV)

	srv.logger.Debugf("Проверка наличия имени метрики...")
	if mN == "" {
		srv.logger.Errorf("Имя метрики не задано")
		srv.logger.Errorf("Отправлен код: %v", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	srv.logger.Debugf("Проверка наличия имени метрики прошла успешно")

	srv.logger.Debugf("Проверка значения метрики...")
	switch mT {
	case constants.GaugeType:
		mVParse, err = strconv.ParseFloat(mV, 64)
		if err != nil || mVParse.(float64) < 0 {
			srv.logger.Errorf("Значение метрики не соответствует требуемому типу float64 или меньше нуля: %v", mV)
			srv.logger.Errorf("Отправлен код: %v", http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case constants.CounterType:
		mVParse, err = strconv.ParseInt(mV, 10, 64)
		if err != nil || mVParse.(int64) < 0 {
			srv.logger.Errorf("Значение метрики не соответствует требуемому типу int64 или меньше нуля: %v", mV)
			srv.logger.Errorf("Отправлен код: %v", http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	default:
		srv.logger.Errorf("Неизвестный тип метрики")
		srv.logger.Errorf("Отправлен код: %v", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	srv.logger.Debugf("Проверка значения метрики прошла успешно")

	receivedMetric := model.Metric{MetType: mT, MetName: mN, MetValue: mVParse}
	srv.logger.Debugf("Полученная метрика: %v", receivedMetric)

	srv.logger.Debugf("Обновление метрики...")
	calcMetric := service.ProcMetric(receivedMetric, store)
	srv.logger.Debugf("Обновление метрики прошло успешно: %v", calcMetric)

	srv.logger.Debugf("Отправлен Content-Type: %v", constants.ContentTypeTPUTF8)
	w.Header().Set("Content-Type", constants.ContentTypeTPUTF8)

	srv.logger.Debugf("Отправлен код: %v", http.StatusOK)
	w.WriteHeader(http.StatusOK)

	srv.logger.Debugf("Обновление хранилища...")
	storage.SM(store, *calcMetric)
	srv.logger.Debugf("Обновление хранилища прошло успешно")
	srv.logger.Debugf("Обновлённое хранилище: %v", storage.GMs(store))
}
