package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golovanevvs/metalecoll/internal/agent/storage/mapstorage"
	"github.com/golovanevvs/metalecoll/internal/server/constants"
)

type agent struct {
	store          mapstorage.Storage
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

			fmt.Println("Получение данных из хранилища...")
			mapStore, err := ag.store.GetMetricsMap()
			if err != nil {
				fmt.Println("Ошибка получения данных из хранилища:", err)
				continue
			}
			fmt.Println(mapStore)
			fmt.Println("Получение данных из хранилища прошло успешно")

			putString = fmt.Sprintf("http://%s/updates/", config.addr)

			fmt.Println("Формирование dto...")

			fmt.Println("Формирование dto прошло успешно")

			fmt.Println("Формирование среза метрик...")

			metrics := make([]Metrics, 0)

			for _, value := range mapStore {
				switch value.MetType {
				case constants.GaugeType:
					v, _ := value.MetValue.(float64)
					body = Metrics{
						ID:    value.MetName,
						MType: value.MetType,
						Value: &v,
					}
				case constants.CounterType:
					v, _ := value.MetValue.(int64)
					body = Metrics{
						ID:    value.MetName,
						MType: value.MetType,
						Delta: &v,
					}
				}

				metrics = append(metrics, body)
			}
			fmt.Println("Формирование среза метрик прошло успешно")
			fmt.Println(metrics)

			fmt.Println("Кодирование в JSON...")
			enc, err := json.Marshal(metrics)
			if err != nil {
				fmt.Println("Ошибка кодирования:", err)
				continue
			}

			fmt.Println("Сжатие в gzip...")
			gzipWr := gzip.NewWriter(&gzipB)
			_, err = gzipWr.Write(enc)
			if err != nil {
				fmt.Println("Ошибка сжатия в gzip:", err)
				gzipWr.Close()
				continue
			}
			gzipWr.Close()
			fmt.Println("Сжатие в gzip прошло успешно")

			fmt.Println("Формирование запроса POST...")
			request, err := http.NewRequest("POST", putString, &gzipB)
			if err != nil {
				fmt.Println("Ошибка формирования запроса:", err)
			}
			fmt.Println("Формирование запроса POST прошло успешно")

			fmt.Println("Установка заголовков...")
			request.Header.Set("Content-Encoding", "gzip")
			request.Header.Set("Content-Type", "application/json")
			fmt.Println("Установка заголовков прошла успешно")

			fmt.Println("Отправка запроса...")
			response, err := client.Do(request)
			if err != nil {
				fmt.Println("Ошибка отправки запроса:", err)
				continue
			}
			response.Body.Close()
			fmt.Println("Отправка запроса прошла успешно")
			fmt.Println("Reporting completed")

			// switch v := value.MetValue.(type) {
			// case float64:
			// 	body = Metrics{
			// 		ID:    value.MetName,
			// 		MType: value.MetType,
			// 		Value: &v,
			// 	}
			// case int64:
			// 	body = Metrics{
			// 		ID:    value.MetName,
			// 		MType: value.MetType,
			// 		Delta: &v,
			// 	}
			// }
			// fmt.Printf("Формирование метрики прошло успешно. Запрос: %v. Метрика: %v\n", putString, body)

			// fmt.Println("Кодирование в JSON...")
			// enc, err := json.Marshal(body)
			// if err != nil {
			// 	fmt.Println("Ошибка кодирования:", err)
			// 	continue
			// }
			// fmt.Println("Кодирование в JSON прошло успешно")

			// fmt.Println("Сжатие в gzip...")
			// gzipWr := gzip.NewWriter(&gzipB)
			// _, err = gzipWr.Write(enc)
			// if err != nil {
			// 	fmt.Println("Ошибка сжатия в gzip:", err)
			// 	gzipWr.Close()
			// 	continue
			// }
			// gzipWr.Close()
			// fmt.Println("Сжатие в gzip прошло успешно")

			// fmt.Println("Формирование запроса POST...")
			// request, err := http.NewRequest("POST", putString, &gzipB)
			// if err != nil {
			// 	fmt.Println("Ошибка формирования запроса:", err)
			// }
			// fmt.Println("Формирование запроса POST прошло успешно")

			// fmt.Println("Установка заголовков...")
			// request.Header.Set("Content-Encoding", "gzip")
			// request.Header.Set("Content-Type", "application/json")
			// fmt.Println("Установка заголовков прошла успешно")

			// fmt.Println("Отправка запроса...")
			// response, err := client.Do(request)
			// if err != nil {
			// 	fmt.Println("Ошибка отправки запроса:", err)
			// 	continue
			// }
			// response.Body.Close()
			// fmt.Println("Отправка запроса прошла успешно")
			// fmt.Println("Reporting completed")

		}
	}
}

func NewAgent(store mapstorage.Storage, pollInterval, reportInterval int) *agent {
	s := &agent{
		store:          store,
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
	}
	return s
}
