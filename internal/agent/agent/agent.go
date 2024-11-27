package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golovanevvs/metalecoll/internal/agent/model"
	"github.com/golovanevvs/metalecoll/internal/agent/storage/mapstorage"
)

type agent struct {
	store          mapstorage.Storage
	pollInterval   int
	reportInterval int
}

func Start(config *config) {
	// store := mapstorage.NewStorage()

	// ag := NewAgent(store, config.pollInterval, config.reportInterval)

	// // запуск таймера обновления метрик
	// pollIntTime := time.NewTicker(time.Duration(ag.pollInterval) * time.Second)
	// defer pollIntTime.Stop()

	// // запуск таймера отправки метрик
	// reportIntTime := time.NewTicker(time.Duration(ag.reportInterval) * time.Second)
	// defer reportIntTime.Stop()

	// // работа таймеров
	// for {
	// 	select {
	// 	case <-pollIntTime.C:
	// 		regNSave(ag)

	// 	case <-reportIntTime.C:
	// 		fmt.Println("-------------------------------------------------------------------------")
	// 		fmt.Println("Reporting...")

	// 		metrics, err := getMetricsSlice(ag)
	// 		if err != nil {
	// 			fmt.Printf("Ошибка получения метрик: %v", err)
	// 			continue
	// 		}

	// 		splitMetrics := splitMetricsSlice(metrics, config.rateLimit)

	// 		urlString := fmt.Sprintf("http://%s/update/", config.addr)

	// 		sendMetrics(splitMetrics, urlString, config.hashKey, config.rateLimit)

	// 		fmt.Println("Reporting completed")
	// 	}
	// }
	var putString string
	//var body Metrics
	var metricsJSONGZIP bytes.Buffer

	store := mapstorage.NewStorage()

	ag := NewAgent(store, config.pollInterval, config.reportInterval)

	client := &http.Client{}

	pollIntTime := time.NewTicker(time.Duration(ag.pollInterval) * time.Second)
	reportIntTime := time.NewTicker(time.Duration(ag.reportInterval) * time.Second)

	defer pollIntTime.Stop()
	defer reportIntTime.Stop()

	for {
		select {
		case <-pollIntTime.C:
			RegisterMetrics(ag)
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

			fmt.Println("Формирование среза метрик...")

			metrics := make([]model.Metric, 0)

			for _, value := range mapStore {
				// switch value.Type {
				// case constants.GaugeType:
				// 	v, _ := value.Value.(float64)
				// 	body = Metrics{
				// 		ID:    value.Name,
				// 		MType: value.Type,
				// 		Value: &v,
				// 	}
				// case constants.CounterType:
				// 	v, _ := value.Value.(int64)
				// 	body = Metrics{
				// 		ID:    value.Name,
				// 		MType: value.Type,
				// 		Delta: &v,
				// 	}
				// }

				//metrics = append(metrics, body)
				metrics = append(metrics, value)
			}
			fmt.Println("Формирование среза метрик прошло успешно")
			fmt.Println(metrics)

			fmt.Println("Кодирование в JSON...")
			metricsJSON, err := json.Marshal(metrics)
			if err != nil {
				fmt.Println("Ошибка кодирования в JSON:", err)
				continue
			}
			fmt.Println("Кодирование в JSON прошло успешно")

			fmt.Println("Сжатие в gzip...")
			gzipWr := gzip.NewWriter(&metricsJSONGZIP)
			_, err = gzipWr.Write(metricsJSON)
			if err != nil {
				fmt.Println("Ошибка сжатия в gzip:", err)
				gzipWr.Close()
				continue
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
				continue
			}
			response.Body.Close()
			fmt.Println("Отправка запроса прошла успешно")
			fmt.Println("Reporting completed")
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
