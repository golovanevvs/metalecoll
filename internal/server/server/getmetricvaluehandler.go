package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golovanevvs/metalecoll/internal/server/constants"
	"github.com/golovanevvs/metalecoll/internal/server/storage"
)

func GetMetricValueHandler(w http.ResponseWriter, r *http.Request) {

	sbody := strings.Split(r.URL.Path, "/")

	if len(sbody) != 4 {
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
	mT := sbody[2] // Тип метрики
	fmt.Println("Тип метрики:", mT)
	mN := sbody[3] // Имя метрики
	fmt.Println("Имя метрики:", mN)

	switch mT {
	case constants.GaugeType, constants.CounterType:
	default:
		fmt.Println("Неизвестный тип метрики")
		fmt.Println("")
		fmt.Println("Отправлен код:", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fmt.Println("Проверка типа метрики прошла успешно")

	fmt.Println("")
	fmt.Println("Получение данных из хранилища...")
	metric, err := storage.GM(srv.store, mN)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Отправлен код:", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	fmt.Println("Получение данных из хранилища прошло успешно")

	value := fmt.Sprintf("%v", metric.MetValue)

	fmt.Println("")
	fmt.Println("Отправлен код:", http.StatusOK)
	w.WriteHeader(http.StatusOK)

	fmt.Println("")
	fmt.Println("Вывод полученных данных...")
	w.Write([]byte(value))
	fmt.Println("Вывод полученных данных прошёл успешно")
}
