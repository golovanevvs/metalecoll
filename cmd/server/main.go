package main

import (
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"
)

const (
	gaugeType    = "gauge"
	counterType  = "counter"
	addr         = ":8080"
	contentType  = "text/plain"
	updateMethod = "update"
)

type metric struct {
	metType  string
	metName  string
	metValue any
}

type MemStorage struct {
	metrics map[string]metric
}

type Storage interface {
	upStore(met metric)
}

var (
	gaugeMet, counterMet metric
	metStorage           MemStorage
	hcount               int
	mT                   string
)

func main() {
	fmt.Println("Запущен сервер:", addr)
	err := http.ListenAndServe(addr, http.HandlerFunc(handlerf))
	if err != nil {
		fmt.Println("Ошибка сервера")
		panic(err)
	}
	fmt.Println(runtime.Version())
}

func (m *metric) updateMetric(mT, mN string, mV any) {
	m.metType = mT
	m.metName = mN
	switch mT {
	case gaugeType:
		m.metValue = mV.(float64)
	case counterType:
		if m.metValue != nil {
			m.metValue = m.metValue.(int64) + mV.(int64)
		} else {
			m.metValue = mV.(int64)
		}
	}
}

func (store *MemStorage) upStore(met metric) {
	if store.metrics == nil {
		store.metrics = make(map[string]metric)
	}
	store.metrics[met.metType] = met
}

func updateStorage(store *MemStorage) {
	switch mT {
	case gaugeType:
		store.upStore(gaugeMet)
	case counterType:
		store.upStore(counterMet)
	}
}

func handlerf(w http.ResponseWriter, r *http.Request) {
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
	case contentType:
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
	mT = sbody[2] // Тип метрики
	fmt.Println("Тип метрики:", mT)
	mN := sbody[3] // Имя метрики
	fmt.Println("Имя метрики:", mN)
	mV := sbody[4] // Значение метрики
	fmt.Println("Значение метрики:", mV)

	fmt.Println("")
	fmt.Println("Проверка типа метода...")

	switch mM {
	case updateMethod:
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
	case gaugeType:
		mVParse, err = strconv.ParseFloat(mV, 64)
		if err != nil || mVParse.(float64) < 0 {
			fmt.Println("Значение метрики не соответствует требуемому типу float64:", mV)
			fmt.Println("")
			fmt.Println("Отправлен код:", http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case counterType:
		mVParse, err = strconv.ParseInt(mV, 10, 64)
		if err != nil || mVParse.(int64) < 0 {
			fmt.Println("Значение метрики не соответствует требуемому типу int64:", mV)
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

	fmt.Println("")
	fmt.Println("Обновление метрики...")

	switch mT {
	case gaugeType:
		gaugeMet.updateMetric(mT, mN, mVParse)
	case counterType:
		counterMet.updateMetric(mT, mN, mVParse)
	}

	fmt.Println("Обновление метрики прошло успешно")

	fmt.Println("")
	fmt.Println("Отправлен Content-Type: text/plain")
	w.Header().Set("Content-Type", "text/plain")

	fmt.Println("")
	fmt.Println("Отправлен код:", http.StatusOK)
	w.WriteHeader(http.StatusOK)

	fmt.Println("")
	fmt.Println("Обновление хранилища...")

	updateStorage(&metStorage)

	fmt.Println("Обновление хранилища прошло успешно")
	fmt.Println("")
	fmt.Println("Обновлённое хранилище:", metStorage)
}

//TODO Обернуть запрос в структуру

//Для запуска теста iter1

//surface
//metricstest-windows-amd64 -test.v -test.run=^TestIteration1$ -binary-path=C:\\dev\\projects\\yapracticum\\metalecoll\\cmd\\server\\server

//home
//metricstest-windows-amd64 -test.v -test.run=^TestIteration1$ -binary-path=C:\\Golovanev\\Dev\\Projects\\YaPracticum\\metalecoll\\cmd\\server
