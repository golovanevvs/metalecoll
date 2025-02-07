package handler

import (
	"encoding/json"
	"net/http"

	"github.com/golovanevvs/metalecoll/internal/server/constants"
	"github.com/golovanevvs/metalecoll/internal/server/dto"
	"github.com/golovanevvs/metalecoll/internal/server/model"
)

// UpdateMetricsJSON - обновление метрики, полученной в JSON.
func (hd *handler) UpdateMetricsJSON(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	decFromCTX, ok := ctx.Value(constants.DecryptKey).([]byte)
	if !ok {
		hd.lg.Errorf("Ошибка декодирования crypto\n")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var mVParse any
	var req, resp dto.Metrics

	hd.lg.Debugf("Декодирование JSON...")
	if err := json.Unmarshal(decFromCTX, &req); err != nil {
		hd.lg.Errorf("Ошибка декодирования JSON: %v", err)
		hd.lg.Errorf("Отправлен код: %v", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	hd.lg.Debugf("Декодирование JSON прошло успешно")

	hd.lg.Debugf("Проверка наличия имени метрики...")
	if req.ID == "" {
		hd.lg.Errorf("Имя метрики не задано")
		hd.lg.Errorf("Отправлен код: %v", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	hd.lg.Debugf("Проверка наличия имени метрики прошла успешно")

	switch req.MType {
	case constants.GaugeType:
		mVParse = *req.Value
	case constants.CounterType:
		mVParse = *req.Delta
	default:
		hd.lg.Errorf("Неизвестный тип метрики")
		hd.lg.Errorf("Отправлен код: %v", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	receivedMetric := model.Metric{MetType: req.MType, MetName: req.ID, MetValue: mVParse}
	hd.lg.Debugf("Полученная метрика: %v", receivedMetric)

	hd.lg.Debugf("Запуск сервиса обновления метрики...")
	calcMetric := hd.sv.UpdateMetric(receivedMetric)
	hd.lg.Debugf("Обновление метрики прошло успешно, обновлённая метрика: %v", calcMetric)

	updatedMap := hd.sv.GetMetricsFromMap()
	hd.lg.Debugf("Обновлённая мапа: %v", updatedMap)

	hd.lg.Debugf("Формирование тела ответа...")
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
	hd.lg.Debugf("Формирование тела ответа прошло успешно")

	hd.lg.Debugf("Кодирование в JSON...")
	enc, err := json.Marshal(resp)
	if err != nil {
		hd.lg.Errorf("Ошибка кодирования: %v", err)
		return
	}
	hd.lg.Debugf("Кодирование в JSON прошло успешно")

	hd.lg.Debugf("Отправлен Content-Type: %v", constants.ContentTypeAJ)
	w.Header().Set("Content-Type", constants.ContentTypeAJ)

	hd.lg.Debugf("Отправлен код: %v", http.StatusOK)
	w.WriteHeader(http.StatusOK)

	w.Write(enc)
}
