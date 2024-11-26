package agent

import (
	"fmt"
	"time"

	"github.com/golovanevvs/metalecoll/internal/agent/storage/mapstorage"
)

type agent struct {
	store          mapstorage.Storage
	pollInterval   int
	reportInterval int
}

func Start(config *config) {
	store := mapstorage.NewStorage()

	ag := NewAgent(store, config.pollInterval, config.reportInterval)

	// запуск таймера обновления метрик
	pollIntTime := time.NewTicker(time.Duration(ag.pollInterval) * time.Second)
	defer pollIntTime.Stop()

	// запуск таймера отправки метрик
	reportIntTime := time.NewTicker(time.Duration(ag.reportInterval) * time.Second)
	defer reportIntTime.Stop()

	// работа таймеров
	for {
		select {
		case <-pollIntTime.C:
			regNSave(ag)

		case <-reportIntTime.C:
			fmt.Println("-------------------------------------------------------------------------")
			fmt.Println("Reporting...")

			metrics, err := getMetricsSlice(ag)
			if err != nil {
				fmt.Printf("Ошибка получения метрик: %v", err)
				continue
			}

			splitMetrics := splitMetricsSlice(metrics, config.rateLimit)

			//urlString := fmt.Sprintf("http://%s/updates/", config.addr)
			urlString := fmt.Sprintf("%s/updates/", config.addr)

			sendMetrics(splitMetrics, urlString, config.hashKey, config.rateLimit)

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
