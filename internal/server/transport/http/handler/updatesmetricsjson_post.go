package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/golovanevvs/metalecoll/internal/server/constants"
	"github.com/golovanevvs/metalecoll/internal/server/dto"
	"github.com/golovanevvs/metalecoll/internal/server/model"
)

// UpdatesMetricsJSON обновление набора метрик, полученных в JSON.
func (hd *handler) UpdatesMetricsJSON(w http.ResponseWriter, r *http.Request) {
	var mVParse any
	var dtoMetrics []dto.Metrics
	var hash string

	hd.lg.Debugf("Декодирование JSON...")
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&dtoMetrics); err != nil {
		hd.lg.Errorf("Ошибка декодирования JSON: %s", err.Error())
		hd.lg.Errorf("Отправлен код: %v", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	hd.lg.Debugf("Декодирование JSON прошло успешно")
	hd.lg.Debugf("%v", dtoMetrics)

	if hd.hashKey != "" {
		hd.lg.Debugf("Проверка соответствия полученного и вычисленного hash...")
		metricsJSON, err := json.Marshal(dtoMetrics)
		if err != nil {
			hd.lg.Debugf("Ошибка кодирования в JSON: %s", err.Error())
			hd.lg.Errorf("Отправлен код: %v", http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		hash = calcHash(metricsJSON, hd.hashKey)
		if hash != r.Header.Get("HashSHA256") {
			hd.lg.Debugf("Проверка соответствия полученного и вычисленного hash прошла успешно")
			hd.lg.Errorf("Вычисленный hash не соответствует полученному")
			hd.lg.Errorf("Отправлен код: %v", http.StatusBadRequest)
			return
		}
		hd.lg.Debugf("Проверка соответствия полученного и вычисленного hash прошла успешно")
		hd.lg.Debugf("Вычисленный hash соответствует полученному")
	}

	for _, m := range dtoMetrics {
		switch m.MType {
		case constants.GaugeType:
			mVParse = *m.Value
		case constants.CounterType:
			mVParse = *m.Delta
		default:
			hd.lg.Errorf("Неизвестный тип метрики: %v", m.MType)
			hd.lg.Errorf("Отправлен код: %v", http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		receivedMetric := model.Metric{MetType: m.MType, MetName: m.ID, MetValue: mVParse}
		hd.lg.Debugf("Полученная метрика: %v", receivedMetric)

		hd.lg.Debugf("Обновление метрики...")
		calcMetric := hd.sv.UpdateMetric(receivedMetric)
		hd.lg.Debugf("Обновление метрики прошло успешно: %v", calcMetric)
	}

	hd.lg.Debugf("Обновлённое хранилище: %v", hd.sv.GetMetricsFromMap())

	hd.lg.Debugf("Установка заголовков...")
	w.Header().Set("Content-Type", constants.ContentTypeAJ)
	if hd.hashKey != "" {
		w.Header().Set("HashSHA256", hash)
	}
	hd.lg.Debugf("Установка заголовков прошла успешно")

	hd.lg.Debugf("Отправлен код: %v", http.StatusOK)
	w.WriteHeader(http.StatusOK)
}

func calcHash(data []byte, key string) string {
	h := sha256.New()
	h.Write(data)
	h.Write([]byte(key))
	dst := h.Sum(nil)
	return hex.EncodeToString(dst)
}
