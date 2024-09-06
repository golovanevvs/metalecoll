package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golovanevvs/metalecoll/internal/server/constants"
	"github.com/golovanevvs/metalecoll/internal/server/dto"
	"github.com/golovanevvs/metalecoll/internal/server/storage"
)

func GetMetricValueJSONHandler(w http.ResponseWriter, r *http.Request, store storage.Storage) {
	var v float64
	var d int64

	srv.logger.Debugf("Декодирование JSON...")

	req := dto.Metrics{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		srv.logger.Errorf("Ошибка декодирования JSON: %v", err)
		srv.logger.Errorf("Отправлен код: %v", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	srv.logger.Debugf("Декодирование JSON прошло успешно: %v", req)

	srv.logger.Debugf("Получение данных из хранилища по name %v...", req.ID)
	metric, err := storage.GM(store, req.ID)
	if err != nil {
		fmt.Println(err)
		srv.logger.Errorf("Ошибка получения данных из хранилища: %v", err)
		srv.logger.Errorf("Отправлен код: %v", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	srv.logger.Debugf("Получение данных из хранилища прошло успешно: %v", metric)

	resp := dto.Metrics{}
	srv.logger.Debugf("Формирование тела ответа...")
	switch metric.MetType {
	case constants.GaugeType:
		v = metric.MetValue.(float64)
		resp = dto.Metrics{
			ID:    metric.MetName,
			MType: metric.MetType,
			Value: &v,
		}
	case constants.CounterType:
		d = metric.MetValue.(int64)
		srv.logger.Debugf("d: %v", d)
		d1 := fmt.Sprintf("%d", d)
		srv.logger.Debugf("d1: %v", d1)
		d2, _ := strconv.Atoi(d1)
		srv.logger.Debugf("d2: %v", d2)
		d3 := int64(d2)
		srv.logger.Debugf("d3: %v", d3)
		resp = dto.Metrics{
			ID:    metric.MetName,
			MType: metric.MetType,
			Delta: &d3,
		}
	}
	srv.logger.Debugf("Формирование тела ответа прошло успешно: %v", resp)

	srv.logger.Debugf("Отправлен Content-Type: %v", constants.ContentTypeAJ)
	w.Header().Set("Content-Type", constants.ContentTypeAJ)

	srv.logger.Debugf("Кодирование в JSON...")

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
