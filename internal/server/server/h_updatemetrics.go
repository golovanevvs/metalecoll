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

	// srv.logger.Debugf("")
	// srv.logger.Debugf("Проверка метода...")

	// // if r.Method != http.MethodPost {
	// 	srv.logger.Errorf("Недопустимый метод: %s", r.Method)
	// 	srv.logger.Errorf("")
	// 	srv.logger.Errorf("Отправлен код: %v", http.StatusMethodNotAllowed)
	// 	w.WriteHeader(http.StatusMethodNotAllowed)
	// 	return
	// }

	// srv.logger.Debugf("Проверка Content-Type...")
	// cT := r.Header.Get("Content-Type")

	// switch cT {
	// case constants.ContentTypeTP, constants.AContentTypeTP, "":
	// default:
	// 	srv.logger.Errorf("Недопустимый content-type: %v", cT)
	// 	srv.logger.Errorf("Отправлен код: %v", http.StatusBadRequest)
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	// srv.logger.Debugf("Проверка Content-Type прошла успешно")

	// srv.logger.Debugf("")
	// srv.logger.Debugf("Чтение и разделение тела запроса...")

	// sbody := strings.Split(r.URL.Path, "/")
	// if len(sbody) != 5 {
	// 	srv.logger.Errorf("Структура тела запроса не соответствует ожидаемой. Получено тело запроса: %v", r.URL.Path)
	// 	srv.logger.Errorf("")
	// 	srv.logger.Errorf("Отправлен код: %v", http.StatusNotFound)
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	// srv.logger.Debugf("Чтение и разделение тела запроса прошло успешно")
	srv.logger.Debugf("Тело запроса: %v", r.URL.Path)

	srv.logger.Debugf("Параметры полученной метрики:")
	// mM := sbody[1] // Тип метода
	// srv.logger.Debugf("Тип метода: %v", mM)
	//mT := sbody[2] // Тип метрики
	mT := chi.URLParam(r, constants.MetTypeURL)
	srv.logger.Debugf("Тип метрики: %v", mT)
	//mN := sbody[3] // Имя метрики
	mN := chi.URLParam(r, constants.MetNameURL)
	srv.logger.Debugf("Имя метрики: %v", mN)
	//mV := sbody[4] // Значение метрики
	mV := chi.URLParam(r, constants.MetValueURL)
	srv.logger.Debugf("Значение метрики: %v", mV)

	// srv.logger.Debugf("")
	// srv.logger.Debugf("Проверка типа метода...")

	// switch mM {
	// case constants.UpdateMethod:
	// default:
	// 	srv.logger.Errorf("Неизвестный тип метода: %v", mM)
	// 	srv.logger.Errorf("")
	// 	srv.logger.Errorf("Отправлен код: %v", http.StatusBadRequest)
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	// srv.logger.Debugf("Проверка типа метода прошла успешно")

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

	srv.logger.Debugf("Обновление метрики...")

	calcMetric := service.ProcMetric(receivedMetric, store)
	srv.logger.Debugf("%v", calcMetric)
	srv.logger.Debugf("Обновление метрики прошло успешно")

	srv.logger.Debugf("Отправлен Content-Type: text/plain; charset=utf-8")
	w.Header().Set("Content-Type", constants.ContentTypeTP)

	srv.logger.Debugf("Отправлен код: %v", http.StatusOK)
	w.WriteHeader(http.StatusOK)

	srv.logger.Debugf("Обновление хранилища...")
	storage.SM(store, *calcMetric)

	srv.logger.Debugf("Обновление хранилища прошло успешно")
	srv.logger.Debugf("Обновлённое хранилище: %v", storage.GMs(store))
}
