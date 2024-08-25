package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/golovanevvs/metalecoll/internal/server/constants"
	"github.com/golovanevvs/metalecoll/internal/server/model"
	"github.com/golovanevvs/metalecoll/internal/server/service"
	"github.com/golovanevvs/metalecoll/internal/server/storage"
)

var (
	hcount int
)

func UpdateMetricsHandler(w http.ResponseWriter, r *http.Request, store storage.Storage) {
	var mVParse any
	var err error

	fmt.Println("")
	fmt.Println("-----------------------------------------------------")
	fmt.Println("")
	hcount++
	fmt.Println("Запрос №", hcount)

	fmt.Println("")
	fmt.Println("Проверка метода...")

	if r.Method != http.MethodPost {
		fmt.Println("Недопустимый метод:", r.Method)
		fmt.Println("")
		fmt.Println("Отправлен код:", http.StatusMethodNotAllowed)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	fmt.Println("Метод:", r.Method)

	fmt.Println("")
	fmt.Println("Проверка Content-Type...")
	cT := r.Header.Get("Content-Type")

	switch cT {
	case constants.ContentType, constants.AContentType, "":
	default:
		fmt.Println("Недопустимый content-type:", cT)
		fmt.Println("")
		fmt.Println("Отправлен код:", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println("Проверка Content-Type прошла успешно")

	fmt.Println("")
	fmt.Println("Чтение и разделение тела запроса...")
	sbody := strings.Split(r.URL.Path, "/")
	if len(sbody) != 5 {
		fmt.Println("Структура тела запроса не соответствует ожидаемой. Получено тело запроса:", r.URL.Path)
		fmt.Println("")
		fmt.Println("Отправлен код:", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fmt.Println("Чтение и разделение тела запроса прошло успешно")
	fmt.Println("")
	fmt.Println("Тело запроса:", r.URL.Path)

	fmt.Println("")
	fmt.Println("Параметры полученной метрики:")
	mM := sbody[1] // Тип метода
	fmt.Println("Тип метода:", mM)
	//mT := sbody[2] // Тип метрики
	mT := chi.URLParam(r, constants.MetTypeURL)
	fmt.Println("Тип метрики:", mT)
	//mN := sbody[3] // Имя метрики
	mN := chi.URLParam(r, constants.MetNameURL)
	fmt.Println("Имя метрики:", mN)
	//mV := sbody[4] // Значение метрики
	mV := chi.URLParam(r, constants.MetValueURL)
	fmt.Println("Значение метрики:", mV)

	fmt.Println("")
	fmt.Println("Проверка типа метода...")

	switch mM {
	case constants.UpdateMethod:
	default:
		fmt.Println("Неизвестный тип метода:", mM)
		fmt.Println("")
		fmt.Println("Отправлен код:", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println("Проверка типа метода прошла успешно")

	fmt.Println("")
	fmt.Println("Проверка наличия имени метрики...")

	if mN == "" {
		fmt.Println("Имя метрики не задано")
		fmt.Println("")
		fmt.Println("Отправлен код:", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fmt.Println("Проверка наличия имени метрики прошла успешно")

	fmt.Println("")
	fmt.Println("Проверка значения метрики...")

	switch mT {
	case constants.GaugeType:
		mVParse, err = strconv.ParseFloat(mV, 64)
		if err != nil || mVParse.(float64) < 0 {
			fmt.Println("Значение метрики не соответствует требуемому типу float64 или меньше нуля:", mV)
			fmt.Println("")
			fmt.Println("Отправлен код:", http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case constants.CounterType:
		mVParse, err = strconv.ParseInt(mV, 10, 64)
		if err != nil || mVParse.(int64) < 0 {
			fmt.Println("Значение метрики не соответствует требуемому типу int64 или меньше нуля:", mV)
			fmt.Println("")
			fmt.Println("Отправлен код:", http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	default:
		fmt.Println("Неизвестный тип метрики")
		fmt.Println("")
		fmt.Println("Отправлен код:", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println("Проверка значения метрики прошла успешно")

	receivedMetric := model.Metric{MetType: mT, MetName: mN, MetValue: mVParse}

	fmt.Println("")
	fmt.Println("Обновление метрики...")

	calcMetric := service.ProcMetric(receivedMetric, store)
	fmt.Println(calcMetric)
	fmt.Println("Обновление метрики прошло успешно")

	fmt.Println("")
	fmt.Println("Отправлен Content-Type: text/plain; charset=utf-8")
	w.Header().Set("Content-Type", constants.ContentType)

	fmt.Println("")
	fmt.Println("Отправлен код:", http.StatusOK)
	w.WriteHeader(http.StatusOK)

	fmt.Println("")
	fmt.Println("Обновление хранилища...")
	storage.SM(store, *calcMetric)

	fmt.Println("Обновление хранилища прошло успешно")
	fmt.Println("")
	fmt.Println("Обновлённое хранилище:", storage.GMs(store))
}
