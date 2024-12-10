package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

// GetMetricValue возвращает значение метрики по её имени.
func (hd *handler) GetMetricValue(w http.ResponseWriter, r *http.Request) {
	hd.lg.Debugf("Тело запроса: %v", r.URL.Path)

	hd.lg.Debugf("Параметры полученной метрики:")
	mT := chi.URLParam(r, "type")
	hd.lg.Debugf("Тип метрики: %v", mT)
	mN := chi.URLParam(r, "name")
	hd.lg.Debugf("Имя метрики: %v", mN)

	hd.lg.Debugf("Получение данных из мапы...")
	metric, err := hd.sv.GetMetricFromMap(mN)
	if err != nil {
		fmt.Println(err)
		hd.lg.Errorf("Ошибка получения данных из мапы")
		hd.lg.Errorf("Отправлен код: %v", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	hd.lg.Debugf("Получение данных из мапы прошло успешно")

	value := fmt.Sprintf("%v", metric.MetValue)

	hd.lg.Debugf("Отправлен код: %v", http.StatusOK)
	w.WriteHeader(http.StatusOK)

	hd.lg.Debugf("Вывод полученных данных...")
	w.Write([]byte(value))
	hd.lg.Debugf("Вывод полученных данных прошёл успешно")
}
