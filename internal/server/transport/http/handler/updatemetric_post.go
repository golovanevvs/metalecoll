package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/golovanevvs/metalecoll/internal/server/constants"
	"github.com/golovanevvs/metalecoll/internal/server/model"
)

// UpdateMetric - обновление метрики.
func (hd *handler) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	var mVParse any
	var err error

	hd.lg.Debugf("Тело запроса: %v", r.URL.Path)

	hd.lg.Debugf("Параметры полученной метрики:")
	mT := chi.URLParam(r, "type")
	hd.lg.Debugf("Тип метрики: %v", mT)
	mN := chi.URLParam(r, "name")
	hd.lg.Debugf("Имя метрики: %v", mN)
	mV := chi.URLParam(r, "value")
	hd.lg.Debugf("Значение метрики: %v", mV)

	hd.lg.Debugf("Проверка наличия имени метрики...")
	if mN == "" {
		hd.lg.Errorf("Имя метрики не задано")
		hd.lg.Errorf("Отправлен код: %v", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	hd.lg.Debugf("Проверка наличия имени метрики прошла успешно")

	hd.lg.Debugf("Проверка значения метрики...")
	switch mT {
	case constants.GaugeType:
		mVParse, err = strconv.ParseFloat(mV, 64)
		if err != nil || mVParse.(float64) < 0 {
			hd.lg.Errorf("Значение метрики не соответствует требуемому типу float64 или меньше нуля: %v", mV)
			hd.lg.Errorf("Отправлен код: %v", http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case constants.CounterType:
		mVParse, err = strconv.ParseInt(mV, 10, 64)
		if err != nil || mVParse.(int64) < 0 {
			hd.lg.Errorf("Значение метрики не соответствует требуемому типу int64 или меньше нуля: %v", mV)
			hd.lg.Errorf("Отправлен код: %v", http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	default:
		hd.lg.Errorf("Неизвестный тип метрики")
		hd.lg.Errorf("Отправлен код: %v", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	hd.lg.Debugf("Проверка значения метрики прошла успешно")

	receivedMetric := model.Metric{MetType: mT, MetName: mN, MetValue: mVParse}
	hd.lg.Debugf("Полученная метрика: %v", receivedMetric)

	hd.lg.Debugf("Запуск сервиса обновления метрики...")
	calcMetric := hd.sv.UpdateMetric(receivedMetric)
	hd.lg.Debugf("Обновление метрики прошло успешно, обновлённая метрика: %v", calcMetric)

	hd.lg.Debugf("Отправлен Content-Type: %v", constants.ContentTypeTPUTF8)
	w.Header().Set("Content-Type", constants.ContentTypeTPUTF8)

	hd.lg.Debugf("Отправлен код: %v", http.StatusOK)
	w.WriteHeader(http.StatusOK)

	upMap := hd.sv.GetMetricsFromMap()
	hd.lg.Debugf("Обновлённая мапа: %v", upMap)
	wr := fmt.Sprintf("type: %s, name: %s, value: %v", upMap[mN].MetType, upMap[mN].MetName, upMap[mN].MetValue)
	w.Write([]byte(wr))

}
