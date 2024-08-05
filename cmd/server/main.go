package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const (
	gaugetype   = "gauge"
	countertype = "counter"
)

type metric struct {
	metType string
	name    string
	value   interface{}
}

type memStorage struct {
	metrics map[string]metric
}

// type Storage interface {
// 	Update(metric Metric) error
// }

type mainHandler struct{}

var (
	gaugeMet, counterMet metric
)

func main() {
	var handler mainHandler
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		panic(err)
	}
}

func (m *metric) update(mT, mN string, mV interface{}) {
	m.metType = mT
	m.name = mN
	switch mT {
	case gaugetype:
		m.value = mV.(float64)
	case countertype:
		if m.value != nil {
			m.value = m.value.(int64) + mV.(int64)
		} else {
			m.value = mV.(int64)
		}
	}
}

func (handle mainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Запущен mainPage")

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	fmt.Println("Чтение Body")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println("Чтение Body прошло успешно")

	sbody := strings.Split(string(body), "/")
	if len(sbody) != 5 {
		fmt.Println("Неправильное тело запроса")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	mT := sbody[2]
	fmt.Println(mT)
	mN := sbody[3]
	if mN == "" {
		fmt.Println("Отсутствует имя метрики")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	fmt.Println(mN)
	mV := sbody[4]
	fmt.Println(mV)

	if mT == gaugetype {
		mVgauge, err := strconv.ParseFloat(mV, 64)
		if err != nil || mVgauge < 0 {
			fmt.Println("Ошибка в ParseFloat")
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		gaugeMet.update(mT, mN, mVgauge)

		fmt.Println(gaugeMet)

	} else if mT == countertype {
		mVcounter, err := strconv.ParseInt(mV, 10, 64)
		if err != nil || mVcounter < 0 {
			fmt.Println("Ошибка в ParseInt")
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		counterMet.update(mT, mN, mVcounter)

		fmt.Println(counterMet)

	} else {
		fmt.Println("Ошибка типа")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println("Статус Ok")
	w.WriteHeader(http.StatusOK)
}
