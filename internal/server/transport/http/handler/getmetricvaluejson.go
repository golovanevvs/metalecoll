package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golovanevvs/metalecoll/internal/server/constants"
	"github.com/golovanevvs/metalecoll/internal/server/dto"
)

func (hd *handler) GetMetricValueJSON(w http.ResponseWriter, r *http.Request) {
	hd.lg.Debugf("Декодирование JSON...")

	req := dto.Metrics{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		hd.lg.Errorf("Ошибка декодирования JSON: %v", err)
		hd.lg.Errorf("Отправлен код: %v", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	hd.lg.Debugf("Декодирование JSON прошло успешно: %v", req)

	hd.lg.Debugf("Получение данных из мапы по name %v...", req.ID)
	metric, err := hd.sv.GetMetricFromMap(req.ID)
	if err != nil {
		fmt.Println(err)
		hd.lg.Errorf("Ошибка получения данных из мапы: %v", err)
		hd.lg.Errorf("Отправлен код: %v", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	hd.lg.Debugf("Получение данных из мапы прошло успешно: %v", metric)

	hd.lg.Debugf("Формирование тела ответа...")
	resp := dto.Metrics{}
	switch metric.MetType {
	case constants.GaugeType:
		v := metric.MetValue.(float64)
		resp = dto.Metrics{
			ID:    metric.MetName,
			MType: metric.MetType,
			Value: &v,
		}
	case constants.CounterType:
		if d, ok := metric.MetValue.(int64); ok {
			hd.lg.Debugf("d: %v", d)
			resp = dto.Metrics{
				ID:    metric.MetName,
				MType: metric.MetType,
				Delta: &d,
			}
		} else {
			d := metric.MetValue.(float64)
			d1 := fmt.Sprintf("%.0f", d)
			d2, _ := strconv.Atoi(d1)
			d3 := int64(d2)
			resp = dto.Metrics{
				ID:    metric.MetName,
				MType: metric.MetType,
				Delta: &d3,
			}
		}
	}
	hd.lg.Debugf("Формирование тела ответа прошло успешно: %v", resp)

	hd.lg.Debugf("Отправлен Content-Type: %v", constants.ContentTypeAJ)
	w.Header().Set("Content-Type", constants.ContentTypeAJ)

	hd.lg.Debugf("Кодирование в JSON...")
	enc, err := json.Marshal(resp)
	if err != nil {
		hd.lg.Errorf("Ошибка кодирования: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	hd.lg.Debugf("Кодирование в JSON прошло успешно")

	hd.lg.Debugf("Отправлен код: %v", http.StatusOK)
	w.WriteHeader(http.StatusOK)
	w.Write(enc)
}
