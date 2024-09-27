package server

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/golovanevvs/metalecoll/internal/server/constants"
	"github.com/golovanevvs/metalecoll/internal/server/dto"
	"github.com/golovanevvs/metalecoll/internal/server/model"
	"github.com/golovanevvs/metalecoll/internal/server/service"
)

func (s server) UpdatesMetricsJSONHandler(w http.ResponseWriter, r *http.Request) {
	var mVParse any
	var req []dto.Metrics

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
	srv.logger.Debugf("%v", req)

	if s.c.Crypto.HashKey != "" {
		srv.logger.Debugf("Проверка hash...")
		body, _ := json.Marshal(req)
		hash := calcHash(body, s.c.Crypto.HashKey)
	}

	for _, m := range req {
		switch m.MType {
		case s.c.MetricTypeNames.GaugeType:
			mVParse = *m.Value
		case s.c.MetricTypeNames.CounterType:
			mVParse = *m.Delta
		default:
			srv.logger.Errorf("Неизвестный тип метрики: %v", m.MType)
			srv.logger.Errorf("Отправлен код: %v", http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		receivedMetric := model.Metric{MetType: m.MType, MetName: m.ID, MetValue: mVParse}
		srv.logger.Debugf("Полученная метрика: %v", receivedMetric)

		srv.logger.Debugf("Обновление метрики...")
		calcMetric := service.ProcMetric(receivedMetric, s.mapStore)
		srv.logger.Debugf("Обновление метрики прошло успешно: %v", calcMetric)

		srv.logger.Debugf("Обновление хранилища...")
		s.mapStore.SaveMetric(*calcMetric)
		srv.logger.Debugf("Обновление хранилища прошло успешно")
	}

	srv.logger.Debugf("Обновлённое хранилище: %v", s.mapStore.GetMetrics())

	srv.logger.Debugf("Отправлен Content-Type: %v", constants.ContentTypeAJ)
	w.Header().Set("Content-Type", constants.ContentTypeAJ)

	srv.logger.Debugf("Отправлен код: %v", http.StatusOK)
	w.WriteHeader(http.StatusOK)
}

func calcHash(data []byte, key string) string {
	h := sha256.New()
	h.Write(data)
	h.Write([]byte(key))
	dst := h.Sum(nil)
	return hex.EncodeToString(dst)
}
