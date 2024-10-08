package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
)

func sendMetWorker(id int, urlString string, hashKey string, metrics <-chan []Metrics, results chan<- string) {
	for m := range metrics {
		var metricsJSONGZIP bytes.Buffer

		client := &http.Client{}

		fmt.Printf("sendMetWorker %v: Кодирование в JSON...\n", id)
		metricsJSON, err := json.Marshal(m)
		if err != nil {
			fmt.Printf("sendMetWorker %v: Ошибка кодирования в JSON: %v\n", id, err)
		}
		fmt.Printf("sendMetWorker %v: Кодирование в JSON прошло успешно\n", id)

		fmt.Printf("sendMetWorker %v: Сжатие в gzip...\n", id)
		gzipWr := gzip.NewWriter(&metricsJSONGZIP)
		_, err = gzipWr.Write(metricsJSON)
		if err != nil {
			fmt.Printf("sendMetWorker %v: Ошибка сжатия в gzip: %v\n", id, err)
			gzipWr.Close()
		}
		gzipWr.Close()
		fmt.Printf("sendMetWorker %v: Сжатие в gzip прошло успешно\n", id)

		fmt.Printf("sendMetWorker %v: Формирование запроса POST...\n", id)
		request, err := http.NewRequest("POST", urlString, &metricsJSONGZIP)
		if err != nil {
			fmt.Printf("sendMetWorker %v: Ошибка формирования запроса: %v\n", id, err)
		}
		fmt.Printf("sendMetWorker %v: Формирование запроса POST прошло успешно\n", id)

		fmt.Printf("sendMetWorker %v: Установка заголовков...\n", id)
		request.Header.Set("Content-Encoding", "gzip")
		request.Header.Set("Content-Type", "application/json")
		if hashKey != "" {
			fmt.Printf("sendMetWorker %v: Формирование hash...\n", id)
			hash := calcHash(metricsJSON, hashKey)
			fmt.Printf("sendMetWorker %v: Формирование hash прошло успешно\n", id)
			request.Header.Set("HashSHA256", hash)
		}
		fmt.Printf("sendMetWorker %v: Установка заголовков прошла успешно\n", id)

		fmt.Printf("sendMetWorker %v: Отправка запроса...\n", id)
		response, err := client.Do(request)
		if err != nil {
			fmt.Printf("sendMetWorker %v: Ошибка отправки запроса: %v\n", id, err)
		}
		response.Body.Close()
		fmt.Printf("sendMetWorker %v: Отправка запроса прошла успешно\n", id)
		results <- fmt.Sprintf("sendMetWorker %v: Reporting completed\n", id)
	}
}
