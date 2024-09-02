package agent

import (
	"bytes"
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

	store := mapstorage.NewStorage()

	ag = NewAgent(store, config.pollInterval, config.reportInterval)

	//client := &http.Client{}

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
				return
			}
			// for _, value := range mapstore {
			// 	putString = fmt.Sprintf("http://%s/%s/%s/%s/%v", config.addr, config.updateMethod, value.MetType, value.MetName, value.MetValue)
			// 	request, err := client.Post(putString, constants.AContentTypeTP, nil)
			// 	if err != nil {
			// 		fmt.Println("Ошибка отправки POST-запроса:", err)
			// 	}
			// 	defer request.Body.Close()
			// }

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
				fmt.Println("Формирование запроса прошло успешно")
				fmt.Println("Кодирование в JSON...")
				enc, err := json.Marshal(body)
				if err != nil {
					fmt.Println("Ошибка кодирования:", err)
					continue
				}
				fmt.Println("Кодирование в JSON прошло успешно")
				fmt.Println("Отправка запроса...")
				request, err := http.Post(putString, "application/json", bytes.NewBuffer(enc))
				if err != nil {
					fmt.Println("Ошибка отправки запроса:", err)
					continue
				}
				fmt.Println("Отправка запроса прошло успешно")
				request.Body.Close()
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
