package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golovanevvs/metalecoll/internal/server/storage"
)

type metNameValue struct {
	BMetName  string `json:"name"`
	BMetValue string `json:"value"`
}

func GetMetricNamesHandler(w http.ResponseWriter, r *http.Request) {
	var metrics []metNameValue

	fmt.Println("")
	fmt.Println("Получение всех известных метрик из хранилища...")
	metricsMap := storage.GMs(srv.store)
	fmt.Println("Получение всех известных метрик из хранилища прошло успешно")

	fmt.Println("")
	fmt.Println("Создание среза имя:значение...")
	for k, v := range metricsMap {
		value := fmt.Sprintf("%v", v.MetValue)
		m := metNameValue{
			BMetName:  k,
			BMetValue: value,
		}
		metrics = append(metrics, m)
	}
	fmt.Println("Создание среза имя:значение прошло успешно")
	fmt.Println(metrics)

	fmt.Println("")
	fmt.Println("Сериализация данных в JSON...")
	resp, err := json.Marshal(metrics)
	if err != nil {
		fmt.Println("Ошибка сериализации данных в JSON")
		return
	}
	fmt.Println("Сериализация данных в JSON прошла успешно")

	fmt.Println("")
	fmt.Println("Устанавливаем заголовок Content-Type для передачи информации, кодированной в JSON")
	w.Header().Set("content-type", "application/json")

	fmt.Println("")
	fmt.Println("Отправляем код:", http.StatusOK)
	w.WriteHeader(http.StatusOK)

	fmt.Println("")
	fmt.Println("Отправка тела ответа...")
	w.Write(resp)
	fmt.Println("Отправка тела ответа прошла успешно")
}
