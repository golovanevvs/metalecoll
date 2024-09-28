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
	var dtoMetrics []dto.Metrics
	var hash string

	srv.logger.Debugf("Декодирование JSON...")
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&dtoMetrics); err != nil {
		srv.logger.Errorf("Ошибка декодирования JSON: %v", err)
		srv.logger.Errorf("Отправлен код: %v", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	srv.logger.Debugf("Декодирование JSON прошло успешно")
	srv.logger.Debugf("%v", dtoMetrics)

	if s.c.Crypto.HashKey != "" {
		srv.logger.Debugf("Проверка соответствия полученного и вычисленного hash...")
		metricsJSON, err := json.Marshal(dtoMetrics)
		if err != nil {
			srv.logger.Debugf("Ошибка кодирования в JSON: %v", err)
			srv.logger.Errorf("Отправлен код: %v", http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		hash = calcHash(metricsJSON, s.c.Crypto.HashKey)
		if hash != r.Header.Get("HashSHA256") {
			srv.logger.Debugf("Проверка соответствия полученного и вычисленного hash прошла успешно")
			s.logger.Errorf("Вычисленный hash не соответствует полученному")
			s.logger.Errorf("Отправлен код: %v", http.StatusBadRequest)
			return
		}
		srv.logger.Debugf("Проверка соответствия полученного и вычисленного hash прошла успешно")
		srv.logger.Debugf("Вычисленный hash соответствует полученному")
	}

	for _, m := range dtoMetrics {
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

	s.logger.Debugf("Установка заголовков...")
	w.Header().Set("Content-Type", constants.ContentTypeAJ)
	if s.c.Crypto.HashKey != "" {
		w.Header().Set("HashSHA256", hash)
	}
	s.logger.Debugf("Установка заголовков прошла успешно")

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
