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

type mainHandler struct{}

var (
	gaugeMet, counterMet metric
	metStorage           memStorage
)

func main() {
	var handler mainHandler
	err := http.ListenAndServe(":8080", handler)
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

func (handle mainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Запущен ServeHTTP")

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// Чтение тела запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Разделение тела запроса
	sbody := strings.Split(string(body), "/")
	if len(sbody) != 5 {
		fmt.Println("Неправильное тело запроса")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	mT := sbody[2] // Тип метрики
	mN := sbody[3] // Имя метрики
	mV := sbody[4] // Значение метрики

	switch mT {
	case gaugeType:
		mVParse, err := strconv.ParseFloat(mV, 64)
		if err != nil || mVParse < 0 {
			fmt.Println("Ошибка в ParseFloat")
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		updateMetrics(&gaugeMet, mT, mN, mVParse)
		fmt.Println(gaugeMet)
	case counterType:
		mVParse, err := strconv.ParseInt(mV, 10, 64)
		if err != nil || mVParse < 0 {
			fmt.Println("Ошибка в ParseInt")
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		updateMetrics(&counterMet, mT, mN, mVParse)
		fmt.Println(counterMet)
	default:
		fmt.Println("Ошибка типа")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println(metStorage)
	fmt.Println("Статус Ok")
	w.WriteHeader(http.StatusOK)
}
