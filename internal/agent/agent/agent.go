package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golovanevvs/metalecoll/internal/agent/storage"
	"github.com/golovanevvs/metalecoll/internal/agent/storage/mapstorage"
)

type agent struct {
	store          storage.Storage
	pollInterval   int
	reportInterval int
}

var ag *agent

func Start(config *config) {
	var putString string
	var body Metrics
	var gzipB bytes.Buffer

	store := mapstorage.NewStorage()

	ag = NewAgent(store, config.pollInterval, config.reportInterval)

	client := &http.Client{}

	pollIntTime := time.NewTicker(time.Duration(ag.pollInterval) * time.Second)
	reportIntTime := time.NewTicker(time.Duration(ag.reportInterval) * time.Second)

	defer pollIntTime.Stop()
	defer reportIntTime.Stop()

	for {
		select {
		case <-pollIntTime.C:
			RegisterMetrics()
		case <-reportIntTime.C:
			fmt.Println("-------------------------------------------------------------------------")
			fmt.Println("Reporting...")

			mapstore, err := storage.GMM(ag.store)
			if err != nil {
				fmt.Println("Ошибка получения данных из хранилища:", err)
				continue
			}

			fmt.Println("Формирование запроса...")
			for _, value := range mapstore {
				putString = fmt.Sprintf("http://%s/%s/", config.addr, config.updateMethod)
				switch v := value.MetValue.(type) {
				case float64:
					body = Metrics{
						ID:    value.MetName,
						MType: value.MetType,
						Value: &v,
					}
				case int64:
					body = Metrics{
						ID:    value.MetName,
						MType: value.MetType,
						Delta: &v,
					}
				}
				fmt.Printf("Формирование запроса прошло успешно. Запрос: %v. Метрика: %v\n", putString, body)
				fmt.Println("Кодирование в JSON...")
				enc, err := json.Marshal(body)
				if err != nil {
					fmt.Println("Ошибка кодирования:", err)
					continue
				}
				fmt.Println("Кодирование в JSON прошло успешно")

				fmt.Println("Сжатие...")
				gzipWr := gzip.NewWriter(&gzipB)
				_, err = gzipWr.Write(enc)
				if err != nil {
					fmt.Println("Ошибка сжатия в gzip:", err)
					gzipWr.Close()
					continue
				}
				fmt.Println("Сжатие прошло успешно")
				gzipWr.Close()

				request, err := http.NewRequest("POST", putString, &gzipB)
				if err != nil {
					fmt.Println("Ошибка формирования запроса:", err)
				}

				request.Header.Set("Content-Encoding", "gzip")
				request.Header.Set("Content-Type", "application/json")

				fmt.Println("Отправка запроса...")
				response, err := client.Do(request)
				if err != nil {
					fmt.Println("Ошибка отправки запроса:", err)
					continue
				}
				fmt.Println("Отправка запроса прошла успешно")
				response.Body.Close()

			}
		}
	}
}

func NewAgent(store storage.Storage, pollInterval, reportInterval int) *agent {
	s := &agent{
		store:          store,
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
	}
	return s
}
