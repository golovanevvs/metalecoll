package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golovanevvs/metalecoll/internal/server/constants"
)

func sendMetrics(ag *agent, config *config) error {
	var putString string
	var body Metrics
	var metricsJSONGZIP bytes.Buffer

	fmt.Println("-------------------------------------------------------------------------")
	fmt.Println("Reporting...")

	client := &http.Client{}

	fmt.Println("Получение данных из хранилища...")
	mapStore, err := ag.store.GetMetricsMap()
	if err != nil {
		fmt.Println("Ошибка получения данных из хранилища:", err)
		return err
	}
	fmt.Println(mapStore)
	fmt.Println("Получение данных из хранилища прошло успешно")

	putString = fmt.Sprintf("http://%s/updates/", config.addr)

	fmt.Println("Формирование среза метрик...")

	metrics := make([]Metrics, 0)

	for _, value := range mapStore {
		switch value.Type {
		case constants.GaugeType:
			v, _ := value.Value.(float64)
			body = Metrics{
				ID:    value.Name,
				MType: value.Type,
				Value: &v,
			}
		case constants.CounterType:
			v, _ := value.Value.(int64)
			body = Metrics{
				ID:    value.Name,
				MType: value.Type,
				Delta: &v,
			}
		}

		metrics = append(metrics, body)
	}
	fmt.Println("Формирование среза метрик прошло успешно")
	fmt.Println(metrics)

	fmt.Println("Кодирование в JSON...")
	metricsJSON, err := json.Marshal(metrics)
	if err != nil {
		fmt.Println("Ошибка кодирования в JSON:", err)
		return err
	}
	fmt.Println("Кодирование в JSON прошло успешно")

	fmt.Println("Сжатие в gzip...")
	gzipWr := gzip.NewWriter(&metricsJSONGZIP)
	_, err = gzipWr.Write(metricsJSON)
	if err != nil {
		fmt.Println("Ошибка сжатия в gzip:", err)
		gzipWr.Close()
		return err
	}
	gzipWr.Close()
	fmt.Println("Сжатие в gzip прошло успешно")

	fmt.Println("Формирование запроса POST...")
	request, err := http.NewRequest("POST", putString, &metricsJSONGZIP)
	if err != nil {
		fmt.Println("Ошибка формирования запроса:", err)
	}
	fmt.Println("Формирование запроса POST прошло успешно")

	fmt.Println("Установка заголовков...")
	request.Header.Set("Content-Encoding", "gzip")
	request.Header.Set("Content-Type", "application/json")
	if config.hashKey != "" {
		fmt.Println("Формирование hash...")
		hash := calcHash(metricsJSON, config.hashKey)
		fmt.Println("Формирование hash прошло успешно")
		request.Header.Set("HashSHA256", hash)
	}
	fmt.Println("Установка заголовков прошла успешно")

	fmt.Println("Отправка запроса...")
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Ошибка отправки запроса:", err)
		return err
	}
	response.Body.Close()
	fmt.Println("Отправка запроса прошла успешно")
	fmt.Println("Reporting completed")
	return nil
}
