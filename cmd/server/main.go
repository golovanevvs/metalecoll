package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	gaugeType   = "gauge"
	counterType = "counter"
	addr        = ":8080"
)

type metric struct {
	metType  string
	metName  string
	metValue interface{}
}

type memStorage struct {
	metrics map[string]metric
}

type MetricsInt interface {
	updateMetric(mT, mN string, mV interface{})
}

var (
	gaugeMet, counterMet metric
	metStorage           memStorage
)

// main
func main() {
	//fmt.Println("Запущен сервер:", addr)
	//fmt.Println("")
	err := http.ListenAndServe(addr, http.HandlerFunc(handlerf))
	if err != nil {
		panic(err)
	}
}

func (m *metric) updateMetric(mT, mN string, mV interface{}) {
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
	metStorage.updateStorage(mT, *m)
}

func (m *memStorage) updateStorage(key string, met metric) {
	if m.metrics == nil {
		m.metrics = make(map[string]metric)
	}
	m.metrics[key] = met
}

func updateMetrics(m MetricsInt, mT, mN string, mV interface{}) {
	m.updateMetric(mT, mN, mV)
}

func handlerf(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("")
	//fmt.Println("Проверка метода")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	//fmt.Println("Метод POST")
	//fmt.Println("")

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// Чтение и разделение тела запроса
	//fmt.Println("Чтение и разделение тела запроса")
	sbody := strings.Split(r.URL.Path, "/")
	if len(sbody) != 5 {
		// 	//fmt.Println("Структура тела запроса не соответствует ожидаемой. Получено тело запроса:", r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	//fmt.Println("Получено тело запроса:", r.URL.Path)
	//fmt.Println("")

	mT := sbody[2] // Тип метрики
	//fmt.Println("Тип метрики:", mT)
	mN := sbody[3] // Имя метрики
	//fmt.Println("Имя метрики:", mN)
	mV := sbody[4] // Значение метрики
	//fmt.Println("Значение метрики:", mV)
	//fmt.Println("")

	switch mT {
	case gaugeType:
		mVParse, err := strconv.ParseFloat(mV, 64)
		if err != nil || mVParse < 0 {
			// fmt.Println("Значение метрики не соответствует требуемому типу float64:", mV)
			// fmt.Println(err)
			// fmt.Println("")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		updateMetrics(&gaugeMet, mT, mN, mVParse)
	case counterType:
		mVParse, err := strconv.ParseInt(mV, 10, 64)
		if err != nil || mVParse < 0 {
			// fmt.Println("Значение метрики не соответствует требуемому типу int64:", mV)
			// fmt.Println(err)
			// fmt.Println("")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		updateMetrics(&counterMet, mT, mN, mVParse)
	default:
		// fmt.Println("Неизвестный тип метрики")
		// fmt.Println("")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println("Обновлённая мапа:", metStorage)
	// fmt.Println("")
	fmt.Println("Отправлен статус Ok")
	// fmt.Println("")
	w.WriteHeader(http.StatusOK)
}

//Для запуска теста iter1
//metricstest-windows-amd64 -test.v -test.run=^TestIteration1$ -binary-path=C:\dev\projects\yapracticum\metalecoll\cmd\server

//metricstest-windows-amd64 -test.v -test.run=^TestIteration1$ -binary-path=C:\Golovanev\Dev\Projects\YaPracticum\metalecoll\cmd\server
