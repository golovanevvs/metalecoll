package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const (
	gaugeType   = "gauge"
	counterType = "counter"
	addr        = "localhost:8080"
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
	http.HandleFunc("/", handler)
	fmt.Println("Запущен сервер:", addr)
	fmt.Println("")
	err := http.ListenAndServe(addr, nil)
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

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("")
	fmt.Println("Проверка метода")

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	fmt.Println("Метод POST")
	fmt.Println("")

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// Чтение тела запроса
	fmt.Println("Чтение тела запроса")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Ошибка чтения тела запроса")
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Разделение тела запроса
	sbody := strings.Split(string(body), "/")
	if len(sbody) != 4 {
		fmt.Println("Структура тела запроса не соответствует ожидаемой. Получено тело запроса:", string(body))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fmt.Println("Получено тело запроса:", string(body))
	fmt.Println("")

	mT := sbody[1] // Тип метрики
	fmt.Println("Тип метрики:", mT)
	mN := sbody[2] // Имя метрики
	fmt.Println("Имя метрики:", mN)
	mV := sbody[3] // Значение метрики
	fmt.Println("Значение метрики:", mV)
	fmt.Println("")

	switch mT {
	case gaugeType:
		mVParse, err := strconv.ParseFloat(mV, 64)
		if err != nil || mVParse < 0 {
			fmt.Println("Значение метрики не соответствует требуемому типу float64:", mV)
			fmt.Println(err)
			fmt.Println("")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		updateMetrics(&gaugeMet, mT, mN, mVParse)
	case counterType:
		mVParse, err := strconv.ParseInt(mV, 10, 64)
		if err != nil || mVParse < 0 {
			fmt.Println("Значение метрики не соответствует требуемому типу int64:", mV)
			fmt.Println(err)
			fmt.Println("")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		updateMetrics(&counterMet, mT, mN, mVParse)
	default:
		fmt.Println("Неизвестный тип метрики")
		fmt.Println("")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println("Обновлённая мапа:", metStorage)
	fmt.Println("")
	fmt.Println("Отправлен статус Ok")
	fmt.Println("")
	w.WriteHeader(http.StatusOK)
}

//Для запуска теста iter1
//metricstest-windows-amd64 -test.v -test.run=^TestIteration1$ -binary-path=C:\dev\projects\yapracticum\metalecoll\cmd\server
