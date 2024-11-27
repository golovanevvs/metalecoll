package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func sendMetWorker(id int, urlString string, hashKey string, metrics <-chan []Metrics, results chan<- string) {
	for m := range metrics {
		client := &http.Client{}

		fmt.Printf("sendMetWorker %d: Кодирование в JSON...\n", id)
		metricsJSON, err := json.Marshal(m)
		if err != nil {
			fmt.Printf("sendMetWorker %d: Ошибка кодирования в JSON: %s\n", id, err.Error())
		}
		fmt.Printf("sendMetWorker %d: Кодирование в JSON прошло успешно\n", id)

		// fmt.Printf("sendMetWorker %d: Сжатие в gzip...\n", id)
		// var metricsJSONGZIP bytes.Buffer
		// gzipWr := gzip.NewWriter(&metricsJSONGZIP)
		// _, err = gzipWr.Write(metricsJSON)
		// if err != nil {
		// 	fmt.Printf("sendMetWorker %d: Ошибка сжатия в gzip: %s\n", id, err.Error())
		// 	gzipWr.Close()
		// }
		// gzipWr.Close()
		// fmt.Printf("sendMetWorker %d: Сжатие в gzip прошло успешно\n", id)

		fmt.Printf("sendMetWorker %d: Формирование запроса POST...\n", id)
		//request, err := http.NewRequest("POST", urlString, &metricsJSONGZIP)
		request, err := http.NewRequest("POST", urlString, bytes.NewBuffer(metricsJSON))
		if err != nil {
			fmt.Printf("sendMetWorker %d: Ошибка формирования запроса: %s\n", id, err.Error())
		}
		fmt.Printf("sendMetWorker %d: Формирование запроса POST прошло успешно\n", id)

		fmt.Printf("sendMetWorker %d: Установка заголовков...\n", id)
		//request.Header.Set("Content-Encoding", "gzip")
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

		//! проверка тела запроса
		// cr, err := gzip.NewReader(request.Body)
		// if err != nil {
		// 	fmt.Printf("sendMetWorker %d: проверка: ошибка декомпрессии запроса: %s\n", id, err.Error())
		// 	return
		// }
		// defer cr.Close()
		fmt.Printf("sendMetWorker %d: проверка: декодирование JSON...\n", id)
		var dm []Metrics
		//dec := json.NewDecoder(cr)
		dec := json.NewDecoder(request.Body)
		if err := dec.Decode(&dm); err != nil {
			fmt.Printf("sendMetWorker %d: gроверка: ошибка декодирования JSON: %s\n", id, err.Error())
			return
		}
		fmt.Printf("sendMetWorker %d: gроверка: полученный JSON: %v\n", id, dm)

		fmt.Printf("url: %s\n", request.URL.String())

		response, err := client.Do(request)
		if err != nil {
			fmt.Printf("sendMetWorker %d: Ошибка отправки запроса: %s\n", id, err.Error())
			return
		}
		response.Body.Close()
		fmt.Printf("sendMetWorker %d: Отправка запроса прошла успешно\n", id)
		results <- fmt.Sprintf("sendMetWorker %d: Reporting completed\n", id)
	}
}
