package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

func sendMetWorker(id int, urlString string, hashKey string, metrics <-chan []Metrics, results chan<- string) {
	mu := new(sync.Mutex)
	for m := range metrics {
		for _, v := range m {
			client := &http.Client{}

			fmt.Printf("sendMetWorker %d: Кодирование в JSON...\n", id)
			metricsJSON, err := json.Marshal(v)
			if err != nil {
				fmt.Printf("sendMetWorker %d: Ошибка кодирования в JSON: %s\n", id, err.Error())
			}
			fmt.Printf("sendMetWorker %d: Кодирование в JSON прошло успешно\n", id)

			fmt.Printf("sendMetWorker %d: Сжатие в gzip...\n", id)
			var metricsJSONGZIP bytes.Buffer
			gzipWr := gzip.NewWriter(&metricsJSONGZIP)
			_, err = gzipWr.Write(metricsJSON)
			if err != nil {
				fmt.Printf("sendMetWorker %d: Ошибка сжатия в gzip: %s\n", id, err.Error())
				gzipWr.Close()
			}
			gzipWr.Close()
			fmt.Printf("sendMetWorker %d: Сжатие в gzip прошло успешно\n", id)

			fmt.Printf("sendMetWorker %d: Формирование запроса POST...\n", id)
			request, err := http.NewRequest("POST", urlString, &metricsJSONGZIP)
			if err != nil {
				fmt.Printf("sendMetWorker %d: Ошибка формирования запроса: %s\n", id, err.Error())
			}
			fmt.Printf("sendMetWorker %d: Формирование запроса POST прошло успешно\n", id)

			fmt.Printf("sendMetWorker %d: Установка заголовков...\n", id)
			request.Header.Set("Content-Encoding", "gzip")
			request.Header.Set("Content-Type", "application/json")
			if hashKey != "" {
				fmt.Printf("sendMetWorker %d: Формирование hash...\n", id)
				hash := calcHash(metricsJSON, hashKey)
				fmt.Printf("sendMetWorker %d: Формирование hash прошло успешно\n", id)
				request.Header.Set("HashSHA256", hash)
			}
			fmt.Printf("sendMetWorker %d: Установка заголовков прошла успешно\n", id)

			fmt.Printf("sendMetWorker %d: Отправка запроса...\n", id)

			request.Close = true

			mu.Lock()
			response, err := client.Do(request)
			if err != nil {
				fmt.Printf("sendMetWorker %d: Ошибка отправки запроса: %s\n", id, err.Error())
				return
			}
			response.Body.Close()
			mu.Unlock()
			fmt.Printf("sendMetWorker %d: Отправка запроса прошла успешно\n", id)
			results <- fmt.Sprintf("sendMetWorker %d: Reporting completed\n", id)
		}
	}
}
