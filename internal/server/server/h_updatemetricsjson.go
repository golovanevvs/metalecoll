package server

import (
	"encoding/json"
	"net/http"

	"github.com/golovanevvs/metalecoll/internal/server/constants"
	"github.com/golovanevvs/metalecoll/internal/server/dto"
	"github.com/golovanevvs/metalecoll/internal/server/model"
	"github.com/golovanevvs/metalecoll/internal/server/service"
)

func (s *server) UpdateMetricsJSONHandler(w http.ResponseWriter, r *http.Request) {
	var mVParse any
	var req, resp dto.Metrics

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
	srv.logger.Debugf("Полученная метрика: %v", receivedMetric)

	srv.logger.Debugf("Обновление метрики...")
	calcMetric := service.ProcMetric(receivedMetric, s.mapStore)
	srv.logger.Debugf("Обновление метрики прошло успешно: %v", calcMetric)

	srv.logger.Debugf("Обновление хранилища...")
	s.mapStore.SaveMetric(*calcMetric)
	srv.logger.Debugf("Обновление хранилища прошло успешно")
	srv.logger.Debugf("Обновлённое хранилище: %v", s.mapStore.GetMetrics())

	srv.logger.Debugf("Формирование тела ответа...")
	switch req.MType {
	case constants.GaugeType:
		if v, ok := calcMetric.MetValue.(float64); ok {
			resp = dto.Metrics{
				ID:    calcMetric.MetName,
				MType: calcMetric.MetType,
				Value: &v,
			}
		}
	case constants.CounterType:
		if v, ok := calcMetric.MetValue.(int64); ok {
			resp = dto.Metrics{
				ID:    calcMetric.MetName,
				MType: calcMetric.MetType,
				Delta: &v,
			}
		}
	}
	srv.logger.Debugf("Формирование тела ответа прошло успешно")

	srv.logger.Debugf("Кодирование в JSON...")
	enc, err := json.Marshal(resp)
	if err != nil {
		srv.logger.Errorf("Ошибка кодирования: %v", err)
		return
	}
	srv.logger.Debugf("Кодирование в JSON прошло успешно")

	srv.logger.Debugf("Отправлен Content-Type: %v", constants.ContentTypeAJ)
	w.Header().Set("Content-Type", constants.ContentTypeAJ)

	srv.logger.Debugf("Отправлен код: %v", http.StatusOK)
	w.WriteHeader(http.StatusOK)

	w.Write(enc)
}
